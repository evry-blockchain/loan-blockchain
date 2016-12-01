package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanTermProposalTableName = "LoanTermProposals"

//Column names
const LTP_LoanTermProposalIDColName = "LoanTermProposalID"
const LTP_LoanTermIDColName = "LoanTermID"
const LTP_ParagraphNumberColName = "ParagraphNumber"
const LTP_LoanTermProposalTextColName = "LoanTermProposalText"
const LTP_LoanTermProposalExpTimeColName = "LoanTermProposalExpTime"

//Column quantity
const LoanTermProposalTableColsQty = 5

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanTermProposalTable(stub shim.ChaincodeStubInterface) error {
	LTP_ColumnNames := []string{LTP_LoanTermProposalIDColName, LTP_LoanTermIDColName, LTP_ParagraphNumberColName,
		LTP_LoanTermProposalTextColName, LTP_LoanTermProposalExpTimeColName}
	return createTable(stub, LoanTermProposalTableName, LTP_ColumnNames)
}

func addLoanTermProposal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission in addLoanTermProposal: " + errA.Error())
	}

	if len(args) == LoanTermProposalTableColsQty {
		return nil, addRow(stub, LoanTermProposalTableName, args, true)
	}
	if len(args) == LoanTermProposalTableColsQty-1 {
		return nil, addRow(stub, LoanTermProposalTableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + "Expecting " + strconv.Itoa(LoanTermProposalTableColsQty-1) +
		" or " + strconv.Itoa(LoanTermProposalTableColsQty))
}

func getLoanTermProposalQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanTermProposalTableName})
}

func getLoanTermProposalList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanTermProposalTableName})
}

func getLoanTermProposalByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanTermProposalTableName, keyValue)
}

func getLoanTermProposalMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanTermProposalTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanTermProposalMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func updateLoanTermProposal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanTermProposalTableColsQty {
		return nil, errors.New("Incorrect number of arguments in updateLoanTermProposal func. Expecting " + strconv.Itoa(LoanTermProposalTableColsQty))
	}

	tbl, err := stub.GetTable(LoanTermProposalTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateLoanTermProposal: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{LoanTermProposalTableName, args[0], cd.Name, args[i]}) //args[0] is hardcoded as row id
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateLoanTermProposal func: " + err.Error())
		}
	}

	return nil, nil
}
