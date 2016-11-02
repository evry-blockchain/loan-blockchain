package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanReturnsTableName = "LoanReturns"

//Column names
const LRT_LoanReturnIDColName = "LoanReturnID"
const LRT_LoanIDColName = "LoanID"
const LRT_AmountColName = "Amount"
const LRT_ReturnDateColName = "ReturnDate"

var LRT_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanReturnTable(stub shim.ChaincodeStubInterface) error {
	LRT_ColumnNames = []string{LRT_LoanReturnIDColName, LRT_LoanIDColName, LRT_AmountColName, LRT_ReturnDateColName}
	return createTable(stub, LoanReturnsTableName, LRT_ColumnNames)
}

func addLoanReturn(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(LRT_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(LRT_ColumnNames)-1))
	}
	err := addRow(stub, LoanReturnsTableName, args, false)
	return nil, err
}

func getLoanReturnsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanReturnsTableName})
}

func getLoanReturnsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanReturnsTableName})
}
