package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const TransactionsTableName = "Transactions"

//Column names
const T_TransactionIDColName = "TransactionID"
const T_FromAccountIDColName = "FromAccountID"
const T_ToAccountIDColName = "ToAccountID "
const T_DateColName = "Date"
const T_TransactionTypeColName = "TransactionType"
const T_TransactionRelatedEntityIDColName = "TransactionRelatedEntityID"
const T_AmountColName = "Amount"

var T_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateTransactionTable(stub *shim.ChaincodeStub) error {
	T_ColumnNames = []string{T_TransactionIDColName, T_FromAccountIDColName, T_ToAccountIDColName,
		T_DateColName, T_TransactionTypeColName, T_TransactionRelatedEntityIDColName, T_AmountColName}
	return createTable(stub, TransactionsTableName, T_ColumnNames)
}

func addTransaction(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != len(T_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(T_ColumnNames)-1))
	}
	err := addRow(stub, TransactionsTableName, args)
	return nil, err
}

func getTransactionsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return countTableRows(stub, []string{TransactionsTableName})
}

func getTransactionsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{TransactionsTableName})
}
