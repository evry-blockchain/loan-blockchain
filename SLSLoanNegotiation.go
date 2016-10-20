package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanNegotiationsTableName = "LoanNegotiations"

//Column names
const LN_LoanNegotiationIDColName = "LoanNegotiationID"
const LN_LoanInvitationIDColName = "LoanInvitationID"
const LN_ParticipantBankIDColName = "ParticipantBankID"
const LN_AmountColName = "Amount"
const LN_NegotiationStatusColName = "NegotiationStatus"
const LN_ParticipantBankCommentColName = "ParticipantBankComment"
const LN_DateColName = "Date"

//Column quantity
const LoanNegotiationsTableColsQty = 7

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanNegotiationTable(stub shim.ChaincodeStubInterface) error {
	LN_ColumnNames := []string{LN_LoanNegotiationIDColName, LN_LoanInvitationIDColName, LN_ParticipantBankIDColName,
		LN_AmountColName, LN_NegotiationStatusColName, LN_ParticipantBankCommentColName, LN_DateColName}
	return createTable(stub, LoanNegotiationsTableName, LN_ColumnNames)
}

func addLoanNegotiation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != LoanNegotiationsTableColsQty-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(LoanNegotiationsTableColsQty-1))
	}

	//Check if related Loan Invitation exists and Arranger BankId is correct in it
	arrangerBankId, err := getTableColValueByKey(stub, LoanInvitationsTableName, args[0], LI_ArrangerBankIDColName) // 0 is a hardcode position of LN_LoanInvitationIDColName argument. Consider avoid hardcoding in the future.
	if err != nil {
		return nil, errors.New("Error getting related Loan Invitation in addLoanNegotiation func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, arrangerBankId)
	if !check {
		return nil, errors.New("Failed checking security in addLoanNegotiation func or returned false: " + err.Error())
	}

	return nil, addRow(stub, LoanNegotiationsTableName, args)
}

func getLoanNegotiationsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanNegotiationsTableName})
}

func getLoanNegotiationsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanNegotiationsTableName})
}

func updateLoanNegotiationStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	loanNegotiationID, newStatus := args[0], args[1]

	//Check if current user has priviledges to update Loan Negotiation Status
	participantBankId, err := getTableColValueByKey(stub, LoanNegotiationsTableName, loanNegotiationID, LN_ParticipantBankIDColName)
	if err != nil {
		return nil, errors.New("Error getting Participant Bank ID in updateLoanNegotiationStatus func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, participantBankId)
	if !check {
		return nil, errors.New("Failed checking security in updateLoanNegotiationStatus func or returned false: " + err.Error())
	}

	return updateTableField(stub, []string{LoanNegotiationsTableName, loanNegotiationID, LN_NegotiationStatusColName, newStatus})
}

func updateParticipantBankComment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	loanNegotiationID, newComment := args[0], args[1]

	//Check if current user has priviledges to update Loan Negotiation Comment
	participantBankId, err := getTableColValueByKey(stub, LoanNegotiationsTableName, loanNegotiationID, LN_ParticipantBankIDColName)
	if err != nil {
		return nil, errors.New("Error getting Participant Bank ID in updateLoanNegotiationStatus func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, participantBankId)
	if !check {
		return nil, errors.New("Failed checking security in updateParticipantBankComment func or returned false: " + err.Error())
	}

	return updateTableField(stub, []string{LoanNegotiationsTableName, loanNegotiationID, LN_ParticipantBankCommentColName, newComment})
}

func getLoanNegotiationByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanNegotiationsTableName, keyValue)
}

func getLoanNegotiationsMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, LoanNegotiationsTableName)
	if err != nil {
		return nil, errors.New("Error in getLoanNegotiationsMaxKey func: " + err.Error())
	}
	return maxKey, nil
}
