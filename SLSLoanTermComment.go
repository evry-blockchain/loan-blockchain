package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanTermCommentTableName = "LoanTermComments"

//Column names
const LTC_LoanTermCommentIDColName = "LoanTermCommentID"
const LTC_ParentLoanTermCommentIDColName = "ParentLoanTermCommentID"
const LTC_LoanTermIDColName = "LoanTermID"
const LTC_UserIDColName = "UserID"
const LTC_BankIDColName = "BankID"
const LTC_CommentTextColName = "CommentText"
const LTC_LoanTermCommentDateColName = "LoanTermCommentDate"

//Column quantity
const LoanTermCommentTableColsQty = 7

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanTermCommentTable(stub shim.ChaincodeStubInterface) error {
	LTC_ColumnNames := []string{LTC_LoanTermCommentIDColName, LTC_ParentLoanTermCommentIDColName, LTC_LoanTermIDColName,
		LTC_UserIDColName, LTC_BankIDColName, LTC_CommentTextColName, LTC_LoanTermCommentDateColName}
	return createTable(stub, LoanTermCommentTableName, LTC_ColumnNames)
}

func addLoanTermComment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission in addLoanTermComment: " + errA.Error())
	}

	if len(args) == LoanTermCommentTableColsQty {
		return nil, addRow(stub, LoanTermCommentTableName, args, true)
	}
	if len(args) == LoanTermCommentTableColsQty-1 {
		return nil, addRow(stub, LoanTermCommentTableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + " ,expected " + strconv.Itoa(LoanTermCommentTableColsQty-1) +
		" or " + strconv.Itoa(LoanTermCommentTableColsQty))
}

func getLoanTermCommentQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanTermCommentTableName})
}

func getLoanTermCommentList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanTermCommentTableName})
}

func getLoanTermCommentByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanTermCommentTableName, keyValue)
}

func getLoanTermCommentMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanTermCommentTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanTermCommentMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func updateLoanTermComment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanTermCommentTableColsQty {
		return nil, errors.New("Incorrect number of arguments in updateLoanTermComment func. Expecting " + strconv.Itoa(LoanTermCommentTableColsQty))
	}

	tbl, err := stub.GetTable(LoanTermCommentTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateLoanTermComment: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{LoanTermCommentTableName, args[0], cd.Name, args[i]}) //args[0] is hardcoded as row id
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateLoanTermComment func: " + err.Error())
		}
	}

	return nil, nil
}
