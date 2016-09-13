package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const ParticipantsTableName = "Participants"
const ParticipantsQuantityKey = "ParticipantsQuantity"

//Column names
const ParticipantKeyColName = "ParticipantKey"
const ParticipantNameColName = "ParticipantName "
const ParticipantTypeColName = "ParticipantType"

//struct for json (consider avoid it)
type Participant struct {
	ParticipantKey  string
	ParticipantName string
	ParticipantType string
}

var ParticipantsQuantity int

// ============================================================================================================================
//
// ============================================================================================================================

func CreateParticipantTable(stub *shim.ChaincodeStub) error {
	//not to forget delete table is it already exists
	stub.DeleteTable(ParticipantsTableName)
	stub.DelState(ParticipantsQuantityKey)

	cdfs := []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: ParticipantKeyColName, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: ParticipantNameColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: ParticipantTypeColName, Type: shim.ColumnDefinition_STRING, Key: false},
	}

	err := stub.CreateTable(ParticipantsTableName, cdfs)
	if err != nil {
		return errors.New("Failed to add participantss table to state: " + err.Error())
	}

	err2 := stub.PutState(ParticipantsQuantityKey, []byte(strconv.Itoa(0)))
	if err2 != nil {
		return errors.New("Failed to add participants quantity to state: " + err2.Error())
	}

	ParticipantsQuantity = 0

	return nil
}

//1. Administrator: add Participant (Bank or Borrower)
//Two arguments expected:
//Participant Name (string)
//Participant Type (string) BANK, BORROWER
func (t *SimpleChaincode) addParticipant(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	qBytes, _ := stub.GetState(ParticipantsQuantityKey)
	if qBytes == nil {
		return nil, errors.New("Participants quantity entity not found")
	}

	qstr := string(qBytes)
	q, errconv := strconv.Atoi(qstr)
	if errconv != nil {
		return nil, errors.New("Error converting key string to int")
	}
	q++
	qstr = strconv.Itoa(q)

	//Add participant to ledger table
	ok, err := stub.InsertRow(ParticipantsTableName, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: qstr}},
			&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}}},
	})
	if !ok && err == nil {
		return nil, errors.New("Participant " + args[0] + " already assigned")
	}

	//Update total participants quantity
	err2 := stub.PutState(ParticipantsQuantityKey, []byte(qstr))

	if err != nil || err2 != nil {
		return nil, errors.New("Failed to add Participant")
	}

	ParticipantsQuantity = q

	s := "The row has been added to Participants table in ledger: \n" +
		"Key: " + qstr + "\n" +
		"Participant: " + qstr + "\n" +
		"Type: " + args[1]
	fmt.Println(s)
	return nil, nil
}

func (t *SimpleChaincode) getParticipantsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	/*qBytes, err := stub.GetState(ParticipantsQuantityKey)
	if qBytes == nil {
		return "", errors.New("Participants quantity entity not found")
	}

	if err != nil {
		return "", errors.New("Failed to get participants quantity")
	}*/

	return []byte(strconv.Itoa(ParticipantsQuantity)), nil
}

func (t *SimpleChaincode) getParticipantsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var p Participant
	var ps []Participant

	for i := 1; i <= ParticipantsQuantity; i++ {
		var cols []shim.Column
		col := shim.Column{Value: &shim.Column_Int32{Int32: int32(i)}}
		cols = append(cols, col)

		row, _ := stub.GetRow(ParticipantsTableName, cols)

		p.ParticipantKey = row.Columns[0].GetString_()
		p.ParticipantName = row.Columns[1].GetString_()
		p.ParticipantType = row.Columns[2].GetString_()

		ps = append(ps, p)
	}

	res, _ := json.Marshal(ps)

	return res, nil
}
