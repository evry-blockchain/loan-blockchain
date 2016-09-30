package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
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

func CreateParticipantTable(stub shim.ChaincodeStubInterface) error {
	P_ColumnNames = []string{P_ParticipantKeyColName, P_ParticipantNameColName, P_ParticipantTypeColName}
	return createTable(stub, ParticipantsTableName, P_ColumnNames)
}

//1. Administrator: add Participant (Bank or Borrower)
//Two arguments expected:
//Participant Name (string)
//Participant Type (string) BANK, BORROWER, LAYER
func addParticipant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != len(P_ColumnNames)-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(len(P_ColumnNames)-1))
	}
	return nil, addRow(stub, ParticipantsTableName, args)
}

func getParticipantsQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{ParticipantsTableName})
}

func getParticipantsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{ParticipantsTableName})
}

func getParticipantsByType(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	filterValue := args[0]
	return filterTableByValue(stub, []string{ParticipantsTableName, P_ParticipantTypeColName, filterValue})
}
