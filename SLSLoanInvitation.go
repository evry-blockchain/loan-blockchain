package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanInvitationsTableName = "LoanInvitations"

//Column names
const LI_LoanInvitationIDColName = "LoanInvitationID"
const LI_ArrangerBankIDColName = "ArrangerBankID"
const LI_BorrowerIDColName = "BorrowerID"
const LI_LoanRequestIDColName = "LoanRequestID"
const LI_LoanTermColName = "LoanTerm"
const LI_AmountColName = "Amount"
const LI_InterestRateColName = "InterestRate"
const LI_InfoColName = "Info"
const LI_StatusColName = "Status"
const LI_Assets = "Assets"
const LI_Convenants = "Convenants"

//Column quantity
const LoanInvitationsTableColsQty = 11

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanInvitationTable(stub shim.ChaincodeStubInterface) error {
	LI_ColumnNames := []string{
		LI_LoanInvitationIDColName, LI_ArrangerBankIDColName, LI_BorrowerIDColName,
		LI_LoanRequestIDColName, LI_LoanTermColName, LI_AmountColName,
		LI_InterestRateColName, LI_InfoColName, LI_StatusColName, LI_Assets, LI_Convenants}
	return createTable(stub, LoanInvitationsTableName, LI_ColumnNames)
}

func addLoanInvitation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanInvitationsTableColsQty-1 {
		return nil, errors.New("Incorrect number of arguments in addLoanInvitation func. Expecting " + strconv.Itoa(LoanInvitationsTableColsQty-1))
	}

	///////////////////////////Security check////////////////////////////
	check, err := checkRowPermissionsByBankId(stub, args[0])
	if !check {
		return nil, errors.New("Failed checking security in addLoanInvitation func or returned false: " + err.Error())
	}
	/////////////////////////////////////////////////////////////////////

	return nil, addRow(stub, LoanInvitationsTableName, args)
}

func updateLoanInvitation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanInvitationsTableColsQty {
		return nil, errors.New("Incorrect number of arguments in addLoanInvitation func. Expecting " + strconv.Itoa(LoanInvitationsTableColsQty))
	}

	loanInvitationID := args[0]

	///////////////////////////Security check////////////////////////////
	check, err := checkLoanInvitationRowPermissionsByBankId(stub, loanInvitationID)
	if !check {
		return nil, errors.New("Failed checking security in updateLoanInvitation or returned false: " + err.Error())
	}
	/////////////////////////////////////////////////////////////////////

	tbl, err := stub.GetTable(LoanInvitationsTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateLoanInvitation: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{LoanInvitationsTableName, loanInvitationID, cd.Name, args[i]})
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateLoanInvitation func: " + err.Error())
		}
	}

	return nil, nil
}

func getLoanInvitationsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanInvitationsTableName})
}

func getLoanInvitationsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanInvitationsTableName})
}

func updateLoanInvitationStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	loanInvitationID, newStatus := args[0], args[1]

	///////////////////////////Security check////////////////////////////
	check, err := checkLoanInvitationRowPermissionsByBankId(stub, loanInvitationID)
	if !check {
		return nil, errors.New("Failed checking security in updateLoanInvitationStatus or returned false: " + err.Error())
	}
	/////////////////////////////////////////////////////////////////////

	return updateTableField(stub, []string{LoanInvitationsTableName, loanInvitationID, LI_StatusColName, newStatus})
}

func getLoanInvitationByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanInvitationsTableName, keyValue)
}

func getLoanInvitationsMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanInvitationsTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanInvitationsMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func checkLoanInvitationRowPermissionsByBankId(stub shim.ChaincodeStubInterface, loanInvitationID string) (bool, error) {
	//Check if current user has priviledges to update Loan Invitation Status
	arrangerBankId, err := getTableColValueByKey(stub, LoanInvitationsTableName, loanInvitationID, LI_ArrangerBankIDColName)
	if err != nil {
		return false, errors.New("Error getting Participant Bank ID in updateLoanNegotiationStatus func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, arrangerBankId)
	if !check {
		return false, errors.New("Failed checking security in updateLoanInvitationStatus func or returned false: " + err.Error())
	}

	return true, nil
}
