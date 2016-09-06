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
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//not to forget delete table is it already exists
	stub.DeleteTable("Participants")

	err := stub.CreateTable("Participants", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Id", Type: shim.ColumnDefinition_UINT32, Key: true},
		&shim.ColumnDefinition{Name: "Name", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_BYTES, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating Participants table: " + err.Error())
	}

	err2 := stub.PutState("ParticipantsQuantity", []byte(strconv.Itoa(0)))
	if err2 != nil {
		return nil, errors.New("Failed to add participants quantity to state")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	if function == "addParticipant" { //initialize the chaincode state, used as reset
		return nil, t.addParticipant(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getParticipantsQuantity" { //read a variable
		q, err := t.getParticipantsQuantity(stub, args)
		if err == nil {
			fmt.Println("Quantity of participants: " + q)
			return nil, nil
		}
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

//1. Administrator: add Participant (Bank or Borrower)
//Two arguments expected:
//Participant Name (string)
//Participant Type (string) BANK, BORROWER
func (t *SimpleChaincode) addParticipant(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 2 {
		return errors.New("Incorrect number of arguments. Expecting 2")
	}

	qBytes, _ := stub.GetState("ParticipantsQuantity")
	if qBytes == nil {
		return errors.New("Participants quantity entity not found")
	}
	q, _ := strconv.Atoi(string(qBytes))
	q++

	//Add participant to ledger table
	ok, err := stub.InsertRow("Participants", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(q)}},
			&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}}},
	})
	if !ok && err == nil {
		return errors.New("Participant " + args[0] + " already assigned")
	}

	//Update total participants quantity
	err2 := stub.PutState("ParticipantsQuantity", []byte(strconv.Itoa(q)))

	if err2 != nil {
		return errors.New("Failed to add Participant")
	}

	fmt.Println("Participant: " + args[0] + " Type: " + args[1] + " has been put to ledger")
	return nil
}

func (t *SimpleChaincode) getParticipantsQuantity(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	qBytes, err := stub.GetState("ParticipantsQuantity")
	if qBytes == nil {
		return "", errors.New("Participants quantity entity not found")
	}

	if err != nil {
		return "", errors.New("Failed to get participants quantity")
	}

	return string(qBytes), nil
}

//2. Arranger Bank: send Loan invitation to Borrower
func (t *SimpleChaincode) sendLoanInvitation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

//3. Borrower: Arranger Bank Loan agree / reject
//4. Arranger Bank: Loan Shares distribution (create many records in Loan Share data structure)/ redistribution (create/edit/remove records in Loan Share data structure)
//5. Participant Bank:  agree/reject Loan Share
//6. Arranger Bank: commit Loan (if all Participant Banks agree) / reject Loan (any time)
//7. Borrower: pay Loan part
//8. Participant Bank: sell Loan entirely / partially (edited)
