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
const LSN_LoanNegotiationIDColName = "LoanNegotiationID"
const LSN_LoanInvitationIDColName = "LoanInvitationID"
const LSN_ParticipantBankIDColName = "ParticipantBankID"
const LSN_AmountColName = "Amount"
const LSN_NegotiationStatusColName = "NegotiationStatus"
const LSN_ParticipantBankCommentColName = "ParticipantBankComment"

var LSN_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanNegotiationTable(stub *shim.ChaincodeStub) error {
	LSN_ColumnNames = []string{LSN_LoanNegotiationIDColName, LSN_LoanInvitationIDColName, LSN_ParticipantBankIDColName,
		LSN_AmountColName, LSN_NegotiationStatusColName, LSN_ParticipantBankCommentColName}
	return createTable(stub, LoanNegotiationsTableName, LSN_ColumnNames)
}

func addLoanNegotiation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != len(LSN_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(LSN_ColumnNames)-1))
	}

	//Check if related Loan Invitation exists and Arranger BankId is correct in it
	arrangerBankId, err := getTableColValueByKey(stub, LoanInvitationsTableName, args[0], LI_ArrangerBankIDColName) // 0 is a hardcode position of LSN_LoanInvitationIDColName argument. Consider avoid hardcoding in the future.
	if err != nil {
		return nil, errors.New("Error getting related Loan Invitation in addLoanNegotiation func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, arrangerBankId)
	if !check {
		return nil, errors.New("Failed checking security in addLoanNegotiation func or returned false: " + err.Error())
	}

	return nil, addRow(stub, LoanNegotiationsTableName, args)
}

func getLoanNegotiationsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanNegotiationsTableName})
}

func getLoanNegotiationsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanNegotiationsTableName})
}

func updateLoanNegotiationStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	loanNegotiationID, newStatus := args[0], args[1]

	//Check if current user has priviledges to update Loan Negotiation Status
	participantBankId, err := getTableColValueByKey(stub, LoanNegotiationsTableName, loanNegotiationID, LSN_ParticipantBankIDColName)
	if err != nil {
		return nil, errors.New("Error getting Participant Bank ID in updateLoanNegotiationStatus func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, participantBankId)
	if !check {
		return nil, errors.New("Failed checking security in updateLoanNegotiationStatus func or returned false: " + err.Error())
	}

	return updateTableField(stub, []string{LoanNegotiationsTableName, loanNegotiationID, LSN_NegotiationStatusColName, newStatus})
}

func updateParticipantBankComment(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	loanNegotiationID, newComment := args[0], args[1]

	//Check if current user has priviledges to update Loan Negotiation Comment
	participantBankId, err := getTableColValueByKey(stub, LoanNegotiationsTableName, loanNegotiationID, LSN_ParticipantBankIDColName)
	if err != nil {
		return nil, errors.New("Error getting Participant Bank ID in updateLoanNegotiationStatus func: " + err.Error())
	}

	check, err := checkRowPermissionsByBankId(stub, participantBankId)
	if !check {
		return nil, errors.New("Failed checking security in updateParticipantBankComment func or returned false: " + err.Error())
	}

	return updateTableField(stub, []string{LoanNegotiationsTableName, loanNegotiationID, LSN_ParticipantBankCommentColName, newComment})
}

func getLoanNegotiationByKey(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, LoanNegotiationsTableName, keyValue)
}
