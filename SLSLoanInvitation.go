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

var LI_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanInvitationTable(stub *shim.ChaincodeStub) error {
	LI_ColumnNames = []string{
		LI_LoanInvitationIDColName, LI_ArrangerBankIDColName, LI_BorrowerIDColName,
		LI_LoanRequestIDColName, LI_LoanTermColName, LI_AmountColName,
		LI_InterestRateColName, LI_InfoColName, LI_StatusColName}
	return createTable(stub, LoanInvitationsTableName, LI_ColumnNames)
}

func addLoanInvitation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != len(LI_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(LI_ColumnNames)-1))
	}
	err := addRow(stub, LoanInvitationsTableName, args)
	return nil, err
}

func getLoanInvitationsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanInvitationsTableName})
}

func getLoanInvitationsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanInvitationsTableName})
}
