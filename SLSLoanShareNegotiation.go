package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanShareNegotiationsTableName = "LoanShareNegotiations"

//Column names
const LSN_LoanShareNegotiationIDColName = "LoanShareNegotiationID"
const LSN_InvitationIDColName = "InvitationID"
const LSN_ParticipantBankIDColName = "ParticipantBankID"
const LSN_AmountColName = "Amount"
const LSN_NegotiationStatusColName = "NegotiationStatus"

var LSN_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanShareNegotiationTable(stub shim.ChaincodeStubInterface) error {
	LSN_ColumnNames = []string{LSN_LoanShareNegotiationIDColName, LSN_InvitationIDColName, LSN_ParticipantBankIDColName,
		LSN_AmountColName, LSN_NegotiationStatusColName}
	return createTable(stub, LoanShareNegotiationsTableName, LSN_ColumnNames)
}

func addLoanShareNegotiation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(LSN_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(LSN_ColumnNames)-1))
	}
	err := addRow(stub, LoanShareNegotiationsTableName, args)
	return nil, err
}

func getLoanShareNegotiationsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{LoanShareNegotiationsTableName})
}

func getLoanShareNegotiationsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{LoanShareNegotiationsTableName})
}
