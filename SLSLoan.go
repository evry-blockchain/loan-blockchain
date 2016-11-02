package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoansTableName = "Loans"

//Column names
const L_LoanIDColName = "LoanID"
const L_ArrangerBankIDColName = "ArrangerBank"
const L_BorrowerIDColName = "BorrowerID"
const L_AmountColName = "Amount"
const L_InterestRateColName = "InterestRate"
const L_InterestTermColName = "InterestTerm"
const L_FeeColName = "Fee"
const L_AgreementDateColName = "AgreementDate"
const L_NegotiationStatusColName = "NegotiationStatus"

var L_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanTable(stub shim.ChaincodeStubInterface) error {
	L_ColumnNames = []string{L_LoanIDColName, L_ArrangerBankIDColName, L_BorrowerIDColName,
		L_AmountColName, L_InterestRateColName, L_InterestTermColName,
		L_FeeColName, L_AgreementDateColName, L_NegotiationStatusColName}
	return createTable(stub, LoansTableName, L_ColumnNames)
}

func addLoan(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(L_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(L_ColumnNames)-1))
	}
	err := addRow(stub, LoansTableName, args, false)
	return nil, err
}

func getLoansQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoansTableName})
}

func getLoansList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoansTableName})
}
