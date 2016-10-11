package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanRequestsTableName = "LoanRequests"

//Column names
const LR_LoanRequestIDColName = "LoanRequestID"
const LR_BorrowerIDColName = "BorrowerID"
const LR_ArrangerBankIDColName = "ArrangerBankID"
const LR_LoanSharesAmountColName = "LoanSharesAmount"
const LR_ProjectRevenueColName = "ProjectRevenue"
const LR_ProjectNameColName = "ProjectName"
const LR_ProjectInformationColName = "ProjectInformation"
const LR_CompanyColName = "Company"
const LR_WebsiteColName = "Website"
const LR_ContactPersonNameColName = "ContactPersonName"
const LR_ContactPersonSurnameColName = "ContactPersonSurname"
const LR_RequestDateColName = "RequestDate"
const LR_StatusColName = "Status"

var LR_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanRequestTable(stub shim.ChaincodeStubInterface) error {
	LR_ColumnNames = []string{LR_LoanRequestIDColName, LR_BorrowerIDColName, LR_ArrangerBankIDColName,
		LR_LoanSharesAmountColName, LR_ProjectRevenueColName, LR_ProjectNameColName, LR_ProjectInformationColName,
		LR_CompanyColName, LR_WebsiteColName, LR_ContactPersonNameColName,
		LR_ContactPersonSurnameColName, LR_RequestDateColName, LR_StatusColName}
	return createTable(stub, LoanRequestsTableName, LR_ColumnNames)
}

func addLoanRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(LR_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments in addLoanRequest func. Expecting " + strconv.Itoa(len(LR_ColumnNames)-1))
	}

	check, err := checkRowPermissionsByBankId(stub, args[1])
	if !check {
		return nil, errors.New("Failed checking security in addLoanRequest func or returned false: " + err.Error())
	}

	return nil, addRow(stub, LoanRequestsTableName, args)
}

func getLoanRequestsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanRequestsTableName})
}

func getLoanRequestsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanRequestsTableName})
}
