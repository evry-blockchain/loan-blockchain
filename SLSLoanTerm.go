package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanTermTableName = "LoanTerms"

//Column names
const LT_LoanTermIDColName = "LoanTermID"
const LT_LoanRequestIDColName = "LoanRequestID"
const LT_ParagraphNumberColName = "ParagraphNumber"
const LT_LoanTermTextColName = "LoanTermText"
const LT_LoanTermStatusColName = "LoanTermStatus"

//Column quantity
const LoanTermTableColsQty = 5

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanTermTable(stub shim.ChaincodeStubInterface) error {
	LT_ColumnNames := []string{LT_LoanTermIDColName, LT_LoanRequestIDColName, LT_ParagraphNumberColName,
		LT_LoanTermTextColName, LT_LoanTermStatusColName}
	return createTable(stub, LoanTermTableName, LT_ColumnNames)
}

func addLoanTerm(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission in addLoanTerm: " + errA.Error())
	}

	if len(args) == LoanTermTableColsQty {
		return nil, addRow(stub, LoanTermTableName, args, true)
	}
	if len(args) == LoanTermTableColsQty-1 {
		return nil, addRow(stub, LoanTermTableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + "Expecting " + strconv.Itoa(LoanTermTableColsQty-1) +
		" or " + strconv.Itoa(LoanTermTableColsQty))
}

func getLoanTermQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanTermTableName})
}

func getLoanTermList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanTermTableName})
}

func getLoanTermByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanTermTableName, keyValue)
}

func getLoanTermMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanTermTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanTermMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func updateLoanTerm(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanTermTableColsQty {
		return nil, errors.New("Incorrect number of arguments in updateLoanTerm func. Expecting " + strconv.Itoa(LoanTermTableColsQty))
	}

	tbl, err := stub.GetTable(LoanTermTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateLoanTerm: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{LoanTermTableName, args[0], cd.Name, args[i]}) //args[0] is hardcoded as row id
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateLoanTerm func: " + err.Error())
		}
	}

	return nil, nil
}
