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
const LR_LoanTermColName = "LoanTerm"
const LR_AssetsColName = "Assets"
const LR_ConvenantsColName = "Convenants"
const LR_InterestRateColName = "InterestRate"

const LoanRequestsTableColsQty = 18

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanRequestTable(stub shim.ChaincodeStubInterface) error {
	LR_ColumnNames := []string{LR_LoanRequestIDColName, LR_BorrowerIDColName, LR_ArrangerBankIDColName,
		LR_LoanSharesAmountColName, LR_ProjectRevenueColName, LR_ProjectNameColName, LR_ProjectInformationColName,
		LR_CompanyColName, LR_WebsiteColName, LR_ContactPersonNameColName,
		LR_ContactPersonSurnameColName, LR_RequestDateColName, LR_StatusColName, LR_MarketAndIndustryColName,
		LR_LoanTermColName, LR_AssetsColName, LR_ConvenantsColName, LR_InterestRateColName}
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

	loanNegStatuses, err := getTableColValuesInSlice(stub, []string{LoanNegotiationsTableName, LN_NegotiationStatusColName, LN_LoanRequestIDColName, loanRequestID})
	if err != nil {
		return errors.New("Error in updateLoanRequestStatus func: " + err.Error())
	}

	var invited, interested, notInterested int
	invited = 0
	interested = 0
	notInterested = 0

	for i, lnStatus := range loanNegStatuses {
		fmt.Println()
		fmt.Println("lnStatus[" + strconv.Itoa(i) + "]: " + lnStatus)
		fmt.Println()
		switch lnStatus {
		case "INTERESTED":
			interested = 1
		case "DECLINED":
			notInterested = 1
		default:
			invited = 1
		}
	}

	// invited | interested | notInterested | statusResult | newLoanRequesStatus
	//    0    |     0      |      0        |     0        |      Draft
	//    0    |     0      |      1        |     1        |      Negotiation Completed
	//    0    |     1      |      0        |     2        |      Negotiation Completed
	//    0    |     1      |      1        |     3        |      Negotiation Completed
	//    1    |     0      |      0        |     4        |      Invitation Sent
	//    1    |     0      |      1        |     5        |      Negotiation Started
	//    1    |     1      |      0        |     6        |      Negotiation Started
	//    1    |     1      |      1        |     7        |      Negotiation Started

	var newLoanRequesStatus string
	statusResult := invited*4 + interested*2 + notInterested
	switch {
	case statusResult == 0:
		newLoanRequesStatus = "Draft"
	case statusResult >= 1 && statusResult <= 3:
		newLoanRequesStatus = "Negotiation Completed"
	case statusResult == 4:
		newLoanRequesStatus = "Invitation Sent"
	case statusResult >= 5 && statusResult <= 7:
		newLoanRequesStatus = "Negotiation Started"
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
