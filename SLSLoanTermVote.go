package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanTermVoteTableName = "LoanTermVotes"

//Column names
const LTV_LoanTermVoteIDColName = "LoanTermVoteID"
const LTV_LoanTermProposalIDColName = "LoanTermProposalID"
const LTV_BankIDColName = "BankID"
const LTV_LoanTermVoteStatusColName = "LoanTermVoteStatus"

//Column quantity
const LoanTermVoteTableColsQty = 4

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanTermVoteTable(stub shim.ChaincodeStubInterface) error {
	LTV_ColumnNames := []string{LTV_LoanTermVoteIDColName, LTV_LoanTermProposalIDColName, LTV_BankIDColName,
		LTV_LoanTermVoteStatusColName}
	return createTable(stub, LoanTermVoteTableName, LTV_ColumnNames)
}

func addLoanTermVote(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission in addLoanTermVote: " + errA.Error())
	}

	if len(args) == LoanTermVoteTableColsQty {
		return nil, addRow(stub, LoanTermVoteTableName, args, true)
	}
	if len(args) == LoanTermVoteTableColsQty-1 {
		return nil, addRow(stub, LoanTermVoteTableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + "Expecting " + strconv.Itoa(LoanTermVoteTableColsQty-1) +
		" or " + strconv.Itoa(LoanTermVoteTableColsQty))
}

func getLoanTermVoteQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanTermVoteTableName})
}

func getLoanTermVoteList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanTermVoteTableName})
}

func getLoanTermVoteByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanTermVoteTableName, keyValue)
}

func getLoanTermVoteMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanTermVoteTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanTermVoteMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func updateLoanTermVote(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanTermVoteTableColsQty {
		return nil, errors.New("Incorrect number of arguments in updateLoanTermVote func. Expecting " + strconv.Itoa(LoanTermVoteTableColsQty))
	}

	tbl, err := stub.GetTable(LoanTermVoteTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateLoanTermVote: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{LoanTermVoteTableName, args[0], cd.Name, args[i]}) //args[0] is hardcoded as row id
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateLoanTermVote func: " + err.Error())
		}
	}

	return nil, nil
}
