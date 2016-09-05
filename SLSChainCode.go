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

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type ParticipantType int

//const (
//	BANK ParticipantType = 1 + iota
//	BORROWER
//)

type Participant struct {
	participantName string
	participantType string
}

var participants []Participant

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
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
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
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" { //read a variable
		fmt.Println("hi there " + function) //error
		return nil, nil
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

//1. Administrator: add Participant (Bank or Borrower)
//Two arguments expected:
//Participant Name (string)
//Participant Type (string) BANK, BORROWER
func (t *SimpleChaincode) addParticipant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var p Participant
	p.participantName = args[0]
	p.participantType = args[1]

	participants = append(participants, p)

	return nil, nil
}

func (t *SimpleChaincode) getParticipantsList(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Incorrect number of arguments. Expecting 2")
	}

	var s string
	for i := 0; i < len(participants); i++ {
		s = s + "Participant Name: " + participants[i].participantName + " Participant Type: " + participants[i].participantType
	}

	return s, nil
}

//2. Arranger Bank: send Loan invitation to Borrower
func (t *SimpleChaincode) sendLoanInvitation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	return nil, nil
}

//3. Borrower: Arranger Bank Loan agree / reject
//4. Arranger Bank: Loan Shares distribution (create many records in Loan Share data structure)/ redistribution (create/edit/remove records in Loan Share data structure)
//5. Participant Bank:  agree/reject Loan Share
//6. Arranger Bank: commit Loan (if all Participant Banks agree) / reject Loan (any time)
//7. Borrower: pay Loan part
//8. Participant Bank: sell Loan entirely / partially (edited)
