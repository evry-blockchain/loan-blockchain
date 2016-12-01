package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanSalesTableName = "LoanSales"

//Column names
const LSL_LoanSaleIDColName = "LoanSaleID"
const LSL_FromLoanShareIDColName = "FromLoanSaleID"
const LSL_ToParticipantIDColName = "ToParticipantID"
const LSL_AmountSoldColName = "AmountSold"

var LSL_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanSaleTable(stub shim.ChaincodeStubInterface) error {
	LSL_ColumnNames = []string{LSL_LoanSaleIDColName, LSL_FromLoanShareIDColName, LSL_ToParticipantIDColName, LSL_AmountSoldColName}
	return createTable(stub, LoanSalesTableName, LSL_ColumnNames)
}

func addLoanSale(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(LSL_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(LSL_ColumnNames)-1))
	}
	err := addRow(stub, LoanSalesTableName, args, false)
	return nil, err
}

func getLoanSalesQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanSalesTableName})
}

func getLoanSalesList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanSalesTableName})
}
