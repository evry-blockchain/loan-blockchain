package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const AccountsTableName = "Accounts"

//Column names
const A_AccountIDColName = "AccountID"
const A_ParticipantIDColName = "ParticipantID"
const A_AmountColName = "Amount"

var A_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func createAccountTable(stub shim.ChaincodeStubInterface) error {
	A_ColumnNames = []string{A_AccountIDColName, A_ParticipantIDColName, A_AmountColName}
	return createTable(stub, AccountsTableName, A_ColumnNames)
}

func addAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(A_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(A_ColumnNames)-1))
	}
	err := addRow(stub, AccountsTableName, args, false)
	return nil, err
}

func getAccountsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{AccountsTableName})
}

func getAccountsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{AccountsTableName})
}

func updateAccountAmount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	accountID, newAmount := args[0], args[1]
	return updateTableField(stub, []string{AccountsTableName, accountID, A_AmountColName, newAmount})
}
