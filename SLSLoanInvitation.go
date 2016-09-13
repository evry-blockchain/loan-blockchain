package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanInvitationsTableName = "LoanInvitations"
const LoanInvitationsQuantityKey = "LoanInvitationsQuantity"

//Column names
const LI_LoanInvitationIDColName = "LoanInvitationID"
const LI_ArrangerBankIDColName = "ArrangerBankID"
const LI_BorrowerIDColName = "BorrowerID"
const LI_LoanRequestIDColName = "LoanRequestID"
const LI_LoanTermColName = "LoanTerm"
const LI_AmountColName = "Amount"
const LI_InterestRateColName = "InterestRate"
const LI_InfoColName = "Info"
const LI_StatusColName = "Status"

const LoanInvitationColumnsQuantity = 9

//struct for json (consider avoid it)
type LoanInvitation struct {
	LoanInvitationID string
	ArrangerBankID   string
	BorrowerID       string
	LoanRequestID    string
	LoanTerm         string
	Amount           string
	InterestRate     string
	Info             string
	Status           string
}

var LoanInvitationsQuantity int

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanInvitationTable(stub *shim.ChaincodeStub) error {
	stub.DeleteTable(LoanInvitationsTableName)

	cdfs := []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: LI_LoanInvitationIDColName, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: LI_ArrangerBankIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_BorrowerIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_LoanRequestIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_LoanTermColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_AmountColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_InterestRateColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_InfoColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LI_StatusColName, Type: shim.ColumnDefinition_STRING, Key: false},
	}

	err := stub.CreateTable(LoanInvitationsTableName, cdfs)
	if err != nil {
		return errors.New("Failed to add Loan Invitations table to state: " + err.Error())
	}

	err2 := stub.PutState(LoanInvitationsQuantityKey, []byte(strconv.Itoa(0)))
	if err2 != nil {
		return errors.New("Failed to add Loan Invitations quantity to state: " + err2.Error())
	}

	LoanInvitationsQuantity = 0

	return nil
}

func (t *SimpleChaincode) addLoanInvitation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != LoanInvitationColumnsQuantity-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(LoanInvitationColumnsQuantity-1))
	}

	qBytes, _ := stub.GetState(LoanInvitationsQuantityKey)
	if qBytes == nil {
		return nil, errors.New("Loan Invitations quantity entity not found")
	}

	qstr := string(qBytes)
	q, errconv := strconv.Atoi(qstr)
	if errconv != nil {
		return nil, errors.New("Error converting key string to int")
	}
	q++
	qstr = strconv.Itoa(q)

	//Add Loan to ledger table
	var cols []*shim.Column
	cols = append(cols, &shim.Column{Value: &shim.Column_String_{String_: qstr}})
	for i := 1; i <= LoanInvitationColumnsQuantity-1; i++ {
		cols = append(cols, &shim.Column{Value: &shim.Column_String_{String_: args[i-1]}})
	}
	ok, err := stub.InsertRow(LoanInvitationsTableName, shim.Row{Columns: cols})
	if !ok && err == nil {
		return nil, errors.New("Loan with LoanID " + qstr + " is already assigned")
	}

	//Update total Loan Invitations quantity
	err2 := stub.PutState(LoanInvitationsQuantityKey, []byte(qstr))

	LoanInvitationsQuantity = q

	if err != nil || err2 != nil {
		return nil, errors.New("Failed to add Loan")
	}

	s := "The row has been added to Loan Invitations table in ledger: \n" +
		LI_LoanInvitationIDColName + " : " + qstr + "\n" +
		LI_ArrangerBankIDColName + " : " + args[0] + "\n" +
		LI_BorrowerIDColName + " : " + args[1] + "\n" +
		LI_LoanRequestIDColName + " : " + args[2] + "\n" +
		LI_LoanTermColName + " : " + args[3] + "\n" +
		LI_AmountColName + " : " + args[4] + "\n" +
		LI_InterestRateColName + " : " + args[5] + "\n" +
		LI_InfoColName + " : " + args[6] + "\n" +
		LI_StatusColName + " : " + args[7]

	fmt.Println(s)
	return nil, nil
}

func (t *SimpleChaincode) getLoanInvitationsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte(strconv.Itoa(LoanInvitationsQuantity)), nil
}

func (t *SimpleChaincode) getLoanInvitationsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var p LoanInvitation
	var ps []LoanInvitation

	for i := 1; i <= LoanInvitationsQuantity; i++ {
		var cols []shim.Column
		col := shim.Column{Value: &shim.Column_Int32{Int32: int32(i)}}
		cols = append(cols, col)

		row, _ := stub.GetRow(LoanInvitationsTableName, cols)

		p.LoanInvitationID = row.Columns[0].GetString_()
		p.ArrangerBankID = row.Columns[1].GetString_()
		p.BorrowerID = row.Columns[2].GetString_()
		p.LoanRequestID = row.Columns[3].GetString_()
		p.LoanTerm = row.Columns[4].GetString_()
		p.Amount = row.Columns[5].GetString_()
		p.InterestRate = row.Columns[6].GetString_()
		p.Info = row.Columns[7].GetString_()
		p.Status = row.Columns[8].GetString_()

		ps = append(ps, p)
	}

	res, _ := json.Marshal(ps)

	return res, nil
}
