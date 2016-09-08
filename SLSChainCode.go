/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const ParticipantsTableName = "Participants"
const ParticipantsQuantityKey = "ParticipantsQuantity"

const ParticipantKeyColName = "ParticipantKey"
const ParticipantNameColName = "ParticipantName "
const ParticipantTypeColName = "ParticipantType"

var ParticipantsQuantity int

type Participant struct {
	ParticipantKey  string `json:"ParticipantKey"`
	ParticipantName string `json:"ParticipantName"`
	ParticipantType string `json:"ParticipantType"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//not to forget delete table is it already exists
	stub.DeleteTable(ParticipantsTableName)
	stub.DelState(ParticipantsQuantityKey)

	err := stub.CreateTable(ParticipantsTableName, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: ParticipantKeyColName, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: ParticipantNameColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: ParticipantTypeColName, Type: shim.ColumnDefinition_STRING, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating Participants table: " + err.Error())
	}

	err2 := stub.PutState(ParticipantsQuantityKey, []byte(strconv.Itoa(0)))
	if err2 != nil {
		return nil, errors.New("Failed to add participants quantity to state")
	}
	ParticipantsQuantity = 0

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	if function == "addParticipant" {
		return t.addParticipant(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getParticipantsQuantity" { //read a variable
		res, err := t.getParticipantsQuantity(stub, args)
		if err != nil {
			return nil, errors.New("Error getting participants quantity")
		}
		return res, nil
	}
	if function == "getParticipantsList" {
		res, err := t.getParticipantsList(stub, args)
		if err != nil {
			return nil, errors.New("Error getting participants list")
		}
		return res, nil
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
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

	s := "The row has been added to Participants table in ledger: Key: " + qstr + " Participant: " + qstr + " Type: " + args[1]
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

//2. Arranger Bank: send Loan invitation to Borrower
func (t *SimpleChaincode) sendLoanInvitation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

//3. Borrower: Arranger Bank Loan agree / reject
//4. Arranger Bank: Loan Shares distribution (create many records in Loan Share data structure)/ redistribution (create/edit/remove records in Loan Share data structure)
//5. Participant Bank:  agree/reject Loan Share
//6. Arranger Bank: commit Loan (if all Participant Banks agree) / reject Loan (any time)
//7. Borrower: pay Loan part
//8. Participant Bank: sell Loan entirely / partially (edited)
