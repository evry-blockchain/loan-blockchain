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

//Column quantity
const ParticipantsTableColsQty = 3

// ============================================================================================================================
//
// ============================================================================================================================

func CreateParticipantTable(stub shim.ChaincodeStubInterface) error {
	P_ColumnNames := []string{P_ParticipantKeyColName, P_ParticipantNameColName, P_ParticipantTypeColName}
	return createTable(stub, ParticipantsTableName, P_ColumnNames)
}

//1. Administrator: add Participant (Bank or Borrower)
//Two arguments expected:
//Participant Name (string)
//Participant Type (string) BANK, BORROWER, LAYER
func addParticipant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission to add Participant: " + errA.Error())
	}

	if len(args) == ParticipantsTableColsQty {
		return nil, addRow(stub, ParticipantsTableName, args, true)
	}
	if len(args) == ParticipantsTableColsQty-1 {
		return nil, addRow(stub, ParticipantsTableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + "Expecting " + strconv.Itoa(ParticipantsTableColsQty-1) +
		" or " + strconv.Itoa(ParticipantsTableColsQty))
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

func getParticipantsByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, ParticipantsTableName, keyValue)
}

func getParticipantsMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, ParticipantsTableName)
	if err != nil {
		return nil, errors.New("Error in getParticipantsMaxKey func: " + err.Error())
	}
	return maxKey, nil
}
