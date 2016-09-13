package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoansTableName = "Loans"
const LoansQuantityKey = "LoansQuantity"

//Column names
const LoanIDColName = "LoanID"
const ArrangerBankIDColName = "ArrangerBank"
const BorrowerIDColName = "BorrowerID"
const AmountColName = "Amount"
const InterestRateColName = "InterestRate"
const InterestTermColName = "InterestTerm"
const FeeColName = "Fee"
const AgreementDateColName = "AgreementDate"
const NegotiationStatusColName = "NegotiationStatus"

//struct for json (consider avoid it)
type Loan struct {
	LoanID            string
	ArrangerBankID    string
	BorrowerID        string
	Amount            string
	InterestRate      string
	InterestTerm      string
	Fee               string
	AgreementDate     string
	NegotiationStatus string
}

var LoansQuantity int

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanTable(stub *shim.ChaincodeStub) error {
	stub.DeleteTable(LoansTableName)

	cdfs := []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: LoanIDColName, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: ArrangerBankIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: BorrowerIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: AmountColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: InterestRateColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: InterestTermColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: FeeColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: AgreementDateColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: NegotiationStatusColName, Type: shim.ColumnDefinition_STRING, Key: false},
	}

	err := stub.CreateTable(LoansTableName, cdfs)
	if err != nil {
		return errors.New("Failed to add loans table to state: " + err.Error())
	}

	err2 := stub.PutState(LoansQuantityKey, []byte(strconv.Itoa(0)))
	if err2 != nil {
		return errors.New("Failed to add loans quantity to state: " + err2.Error())
	}

	LoansQuantity = 0

	return nil
}

func (t *SimpleChaincode) addLoan(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	qBytes, _ := stub.GetState(LoansQuantityKey)
	if qBytes == nil {
		return nil, errors.New("Loans quantity entity not found")
	}

	qstr := string(qBytes)
	q, errconv := strconv.Atoi(qstr)
	if errconv != nil {
		return nil, errors.New("Error converting key string to int")
	}
	q++
	qstr = strconv.Itoa(q)

	//Add Loan to ledger table
	ok, err := stub.InsertRow(LoansTableName, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: qstr}},
			&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[2]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[3]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[4]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[5]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[6]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[7]}}},
	})
	if !ok && err == nil {
		return nil, errors.New("Loan with LoanID " + qstr + " is already assigned")
	}

	//Update total loans quantity
	err2 := stub.PutState(LoansQuantityKey, []byte(qstr))

	if err != nil || err2 != nil {
		return nil, errors.New("Failed to add Loan")
	}

	LoansQuantity = q

	s := "The row has been added to Loans table in ledger: \n" +
		"LoanID: " + qstr + "\n" +
		" ArrangerBankID: " + args[0] + "\n" +
		" BorrowerID: " + args[1] + "\n" +
		" Amount: " + args[2] + "\n" +
		" InterestRate: " + args[3] + "\n" +
		" InterestTerm: " + args[4] + "\n" +
		" Fee: " + args[5] + "\n" +
		" AgreementDate: " + args[6] + "\n" +
		" NegotiationStatus: " + args[7]

	fmt.Println(s)
	return nil, nil
}

func (t *SimpleChaincode) getLoansQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte(strconv.Itoa(LoansQuantity)), nil
}

func (t *SimpleChaincode) getLoansList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var p Loan
	var ps []Loan

	for i := 1; i <= LoansQuantity; i++ {
		var cols []shim.Column
		col := shim.Column{Value: &shim.Column_Int32{Int32: int32(i)}}
		cols = append(cols, col)

		row, _ := stub.GetRow(LoansTableName, cols)

		p.LoanID = row.Columns[0].GetString_()
		p.ArrangerBankID = row.Columns[1].GetString_()
		p.BorrowerID = row.Columns[2].GetString_()
		p.Amount = row.Columns[3].GetString_()
		p.InterestRate = row.Columns[4].GetString_()
		p.InterestTerm = row.Columns[5].GetString_()
		p.Fee = row.Columns[6].GetString_()
		p.AgreementDate = row.Columns[7].GetString_()
		p.NegotiationStatus = row.Columns[8].GetString_()

		ps = append(ps, p)
	}

	res, _ := json.Marshal(ps)

	return res, nil
}
