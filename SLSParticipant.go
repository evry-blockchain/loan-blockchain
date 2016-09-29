package main

import (
	//"encoding/json"
	"errors"
	//	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const ParticipantsTableName = "Participants"

//Column names
const P_ParticipantKeyColName = "ParticipantKey"
const P_ParticipantNameColName = "ParticipantName"
const P_ParticipantTypeColName = "ParticipantType"

var P_ColumnNames []string

// ============================================================================================================================
//
// ============================================================================================================================

func CreateParticipantTable(stub *shim.ChaincodeStub) error {
	P_ColumnNames = []string{P_ParticipantKeyColName, P_ParticipantNameColName, P_ParticipantTypeColName}
	return createTable(stub, ParticipantsTableName, P_ColumnNames)
}

//1. Administrator: add Participant (Bank or Borrower)
//Two arguments expected:
//Participant Name (string)
//Participant Type (string) BANK, BORROWER, LAYER
func addParticipant(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// Verify the identity of the caller
	// Only an administrator can invoker assign
	adminCertificate, err := stub.GetState("admin")
	if err != nil {
		return nil, errors.New("Failed fetching admin identity")
	}

	ok, err := isCaller(stub, adminCertificate)
	if err != nil {
		return nil, errors.New("Failed checking admin identity")
	}
	if !ok {
		return nil, errors.New("The caller is not an administrator")
	}

	if len(args) != len(P_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(P_ColumnNames)-1))
	}
	err = addRow(stub, ParticipantsTableName, args)
	return nil, err
}

func getParticipantsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return countTableRows(stub, []string{ParticipantsTableName})
}

func getParticipantsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{ParticipantsTableName})
}

func getParticipantsByType(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	filterValue := args[0]
	return filterTableByValue(stub, []string{ParticipantsTableName, P_ParticipantTypeColName, filterValue})
}
