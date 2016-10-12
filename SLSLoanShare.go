package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanSharesTableName = "LoanShares"

//Column names

const LS_LoanShareIdColName = "LoanShareId"
const LS_LoanIdColName = "LoanId"
const LS_ParticipantBankIdColName = "ParticipantBankId"
const LS_AmountColName = "Amount"
const LS_NegotiationStatusColName = "NegotiationStatus"

var LS_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanSharesTable(stub *shim.ChaincodeStub) error {
	LS_ColumnNames = []string{LS_LoanShareIdColName, LS_LoanIdColName, LS_ParticipantBankIdColName,
		LS_AmountColName, LS_NegotiationStatusColName}
	return createTable(stub, LoanSharesTableName, LS_ColumnNames)
}

func addLoanShare(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != len(LS_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(LS_ColumnNames)-1))
	}
	err := addRow(stub, LoanSharesTableName, args)
	return nil, err
}

func getLoanSharesQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanSharesTableName})
}

func getLoanSharesList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanSharesTableName})
}
