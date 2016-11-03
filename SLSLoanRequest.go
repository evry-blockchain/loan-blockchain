package main

import (
	//"encoding/json"
	"errors"
	"fmt"
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
const LR_MarketAndIndustryColName = "MarketAndIndustry"

const LoanRequestsTableColsQty = 14

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanRequestTable(stub shim.ChaincodeStubInterface) error {
	LR_ColumnNames := []string{LR_LoanRequestIDColName, LR_BorrowerIDColName, LR_ArrangerBankIDColName,
		LR_LoanSharesAmountColName, LR_ProjectRevenueColName, LR_ProjectNameColName, LR_ProjectInformationColName,
		LR_CompanyColName, LR_WebsiteColName, LR_ContactPersonNameColName,
		LR_ContactPersonSurnameColName, LR_RequestDateColName, LR_StatusColName, LR_MarketAndIndustryColName}
	return createTable(stub, LoanRequestsTableName, LR_ColumnNames)
}

func addLoanRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanRequestsTableColsQty-1 {
		return nil, errors.New("Incorrect number of arguments in addLoanRequest func. Expecting " + strconv.Itoa(LoanRequestsTableColsQty-1))
	}

	///////////////////////////Security check////////////////////////////
	check, err := checkRowPermissionsByBankId(stub, args[1])
	if !check {
		return nil, errors.New("Failed checking security in addLoanRequest func or returned false: " + err.Error())
	}
	/////////////////////////////////////////////////////////////////

	return nil, addRow(stub, LoanRequestsTableName, args, false)
}

func updateLoanRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanRequestsTableColsQty {
		return nil, errors.New("Incorrect number of arguments in updateLoanRequest func. Expecting " + strconv.Itoa(LoanRequestsTableColsQty))
	}

	loanRequestID := args[0]

	///////////////////////////Security check////////////////////////////
	check, err := checkLoanRequestRowPermissionsByBankId(stub, loanRequestID)
	if !check {
		return nil, errors.New("Failed checking security in updateLoanRequest or returned false: " + err.Error())
	}
	/////////////////////////////////////////////////////////////////////

	tbl, err := stub.GetTable(LoanRequestsTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateLoanRequest: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{LoanRequestsTableName, loanRequestID, cd.Name, args[i]})
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateLoanRequest func: " + err.Error())
		}
	}

	return nil, nil
}

func getLoanRequestsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanRequestsTableName})
}

func getLoanRequestsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanRequestsTableName})
}

func getLoanRequestByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanRequestsTableName, keyValue)
}

func getLoanRequestsMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanRequestsTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanRequestsMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func checkLoanRequestRowPermissionsByBankId(stub shim.ChaincodeStubInterface, loanRequestID string) (bool, error) {
	arrangerBankId, err := getTableColValueByKey(stub, LoanRequestsTableName, loanRequestID, LR_ArrangerBankIDColName)
	if err != nil {
		return false, errors.New("Error getting Arranger Bank ID in checkLoanRequestRowPermissionsByBankId func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, arrangerBankId)
	if !check {
		return false, errors.New("Failed checking security in checkLoanRequestRowPermissionsByBankId func or returned false: " + err.Error())
	}

	return true, nil
}

func updateLoanRequestStatus(stub shim.ChaincodeStubInterface, loanRequestID string) error {

	loanInvitationIDs, err := getTableColValuesInSlice(stub, []string{LoanInvitationsTableName, LI_LoanInvitationIDColName, LI_LoanRequestIDColName, loanRequestID})
	if err != nil {
		return errors.New("Error in updateLoanRequestStatus func: " + err.Error())
	}

	var loanNegStatuses []string
	for _, liID := range loanInvitationIDs {
		lnStatuses, _ := getTableColValuesInSlice(stub, []string{LoanNegotiationsTableName, LN_NegotiationStatusColName, LN_LoanInvitationIDColName, liID})
		if err != nil {
			return errors.New("Error in updateLoanRequestStatus func: " + err.Error())
		}
		loanNegStatuses = append(loanNegStatuses, lnStatuses...)
	}

	var invited, interested, notInterested bool
	invited = false
	interested = false
	notInterested = false

	for i, lnStatus := range loanNegStatuses {
		fmt.Println()
		fmt.Println("lnStatus[" + strconv.Itoa(i) + "]: " + lnStatus)
		fmt.Println()
		switch lnStatus {
		case "INTERESTED":
			interested = true
		case "DECLINED":
			notInterested = true
		default:
			invited = true
		}
	}

	var newLoanRequesStatus string
	if invited && !interested && !notInterested {
		newLoanRequesStatus = "Invitation Sent"
	} else {
		if invited && (interested || notInterested) {
			newLoanRequesStatus = "Negotiation Started"
		} else {
			if !invited && (interested || notInterested) {
				newLoanRequesStatus = "Negotiation Completed"
			} else {
				newLoanRequesStatus = "Undefined"
			}
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Println(" newLoanRequesStatus: " + newLoanRequesStatus)
	fmt.Println()
	fmt.Println()

	_, err = updateTableField(stub, []string{LoanRequestsTableName, loanRequestID, LR_StatusColName, newLoanRequesStatus})
	if err != nil {
		return errors.New("Error in updateLoanRequestStatus func: " + err.Error())
	}

	return nil
}

func getTableColValuesInSlice(stub shim.ChaincodeStubInterface, args []string) ([]string, error) {

	// 2 or 4 arguments should be provided:
	// 2 - filter is not needed: tableName, columnName
	// 4 - filter is need: tableName, columnName, filterColumn, filterValue

	var tableName, columnName, filterColumn, filterValue string
	var tbl *shim.Table
	var rows []shim.Row
	var err error

	switch l := len(args); l {
	case 2:
		tableName, filterColumn = args[0], args[1]
		tbl, rows, err = getRowsByColumnValue(stub, []string{tableName})
	case 4:
		tableName, columnName, filterColumn, filterValue = args[0], args[1], args[2], args[3]
		tbl, rows, err = getRowsByColumnValue(stub, []string{tableName, filterColumn, filterValue})
	default:
		return nil, errors.New("Incorrect number of arguments in getTableColValuesInSlice func. Expecting: 2 or 4, " +
			"provided: " + strconv.Itoa(l))
	}

	if err != nil {
		return nil, errors.New("Error in getTableColValuesInSlice func: " + err.Error())
	}

	var colID int
	var isColFound bool
	isColFound = false
	for i, cd := range tbl.ColumnDefinitions {
		if cd.Name == columnName {
			colID = i
			isColFound = true
		}
	}
	if !isColFound {
		return nil, errors.New("Error in getTableColValuesInSlice func: Column '" + columnName + "' is not found in '" + tableName + "' table")
	}

	var colValues []string
	for _, row := range rows {
		colValues = append(colValues, row.Columns[colID].GetString_())
	}
	return colValues, nil
}
