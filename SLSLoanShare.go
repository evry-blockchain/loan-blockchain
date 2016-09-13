package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanSharesTableName = "LoanShares"
const LoanSharesQuantityKey = "LoanSharesQuantity"

//Column names

const LS_LoanShareIdColName = "LoanShareId"
const LS_LoanIdColName = "LoanId"
const LS_ParticipantBankIdColName = "ParticipantBankId"
const LS_AmountColName = "Amount"
const LS_NegotiationStatusColName = "NegotiationStatus"

//struct for json (consider avoid it)
type LoanShare struct {
	LoanShareId       string
	LoanId            string
	ParticipantBankId string
	Amount            string
	NegotiationStatus string
}

var LoanSharesQuantity int

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanSharesTable(stub *shim.ChaincodeStub) error {
	stub.DeleteTable(LoanSharesTableName)

	cdfs := []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: LS_LoanShareIdColName, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: LS_LoanIdColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LS_ParticipantBankIdColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LS_AmountColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LS_NegotiationStatusColName, Type: shim.ColumnDefinition_STRING, Key: false},
	}

	err := stub.CreateTable(LoanSharesTableName, cdfs)
	if err != nil {
		return errors.New("Failed to add LoanShares table to state: " + err.Error())
	}

	err2 := stub.PutState(LoanSharesQuantityKey, []byte(strconv.Itoa(0)))
	if err2 != nil {
		return errors.New("Failed to add LoanShares quantity to state: " + err2.Error())
	}

	LoanSharesQuantity = 0

	return nil
}

func (t *SimpleChaincode) addLoanShare(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	qBytes, _ := stub.GetState(LoanSharesQuantityKey)
	if qBytes == nil {
		return nil, errors.New("LoanShares quantity entity not found")
	}

	qstr := string(qBytes)
	q, errconv := strconv.Atoi(qstr)
	if errconv != nil {
		return nil, errors.New("Error converting key string to int")
	}
	q++
	qstr = strconv.Itoa(q)

	//Add LoanShare to ledger table
	ok, err := stub.InsertRow(LoanSharesTableName, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: qstr}},
			&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[2]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[3]}}},
	})
	if !ok && err == nil {
		return nil, errors.New("LoanShare with LoanShareID " + qstr + " is already assigned")
	}

	//Update total LoanShares quantity
	err2 := stub.PutState(LoanSharesQuantityKey, []byte(qstr))

	if err != nil || err2 != nil {
		return nil, errors.New("Failed to add LoanShare")
	}

	LoanSharesQuantity = q

	fmt.Println("q=" + strconv.Itoa(q))

	s := "The row has been added to LoanShares table in the ledger: \n" +
		"LoanShareId: " + qstr + "\n" +
		" LoanId: " + args[0] + "\n" +
		" ParticipantBankId: " + args[1] + "\n" +
		" Amount: " + args[2] + "\n" +
		" NegotiationStatus: " + args[3]

	fmt.Println(s)
	return nil, nil
}

func (t *SimpleChaincode) getLoanSharesQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte(strconv.Itoa(LoanSharesQuantity)), nil
}

func (t *SimpleChaincode) getLoanSharesList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var p LoanShare
	var ps []LoanShare

	for i := 1; i <= LoanSharesQuantity; i++ {
		var cols []shim.Column
		col := shim.Column{Value: &shim.Column_Int32{Int32: int32(i)}}
		cols = append(cols, col)

		row, _ := stub.GetRow(LoanSharesTableName, cols)

		p.LoanShareId = row.Columns[0].GetString_()
		p.LoanId = row.Columns[1].GetString_()
		p.ParticipantBankId = row.Columns[2].GetString_()
		p.Amount = row.Columns[3].GetString_()
		p.NegotiationStatus = row.Columns[4].GetString_()

		ps = append(ps, p)
	}

	res, _ := json.Marshal(ps)

	return res, nil
}
