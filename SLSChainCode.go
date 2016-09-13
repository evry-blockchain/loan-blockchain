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

	errp := CreateParticipantTable(stub)
	if errp != nil {
		return nil, errors.New("Failed creating Participants table: " + errp.Error())
	}

	errl := CreateLoanTable(stub)
	if errl != nil {
		return nil, errors.New("Failed creating Loan table: " + errl.Error())
	}
	errls := CreateLoanSharesTable(stub)
	if errls != nil {
		return nil, errors.New("Failed creating LoanShares table: " + errls.Error())
	}
	errlr := CreateLoanRequestTable(stub)
	if errlr != nil {
		return nil, errors.New("Failed creating LoanRequests table: " + errlr.Error())
	}
	errli := CreateLoanInvitationTable(stub)
	if errli != nil {
		return nil, errors.New("Failed creating LoanInvitations table: " + errli.Error())
	}

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
	if function == "addLoan" {
		return t.addLoan(stub, args)
	}
	if function == "addLoanShare" {
		return t.addLoanShare(stub, args)
	}
	if function == "addLoanRequest" {
		return t.addLoanRequest(stub, args)
	}
	if function == "addLoanInvitation" {
		return t.addLoanInvitation(stub, args)
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
	if function == "getLoansQuantity" { //read a variable
		res, err := t.getLoansQuantity(stub, args)
		if err != nil {
			return nil, errors.New("Error getting participants quantity")
		}
		return res, nil
	}
	if function == "getLoansList" {
		res, err := t.getLoansList(stub, args)
		if err != nil {
			return nil, errors.New("Error getting participants list")
		}
		return res, nil
	}
	if function == "getLoanSharesQuantity" { //read a variable
		res, err := t.getLoanSharesQuantity(stub, args)
		if err != nil {
			return nil, errors.New("Error getting participants quantity")
		}
		return res, nil
	}
	if function == "getLoanSharesList" {
		res, err := t.getLoanSharesList(stub, args)
		if err != nil {
			return nil, errors.New("Error getting participants list")
		}
		return res, nil
	}
	if function == "getLoanRequestsQuantity" { //read a variable
		res, err := t.getLoanRequestsQuantity(stub, args)
		if err != nil {
			return nil, errors.New("Error getting loan requests quantity")
		}
		return res, nil
	}
	if function == "getLoanRequestsList" {
		res, err := t.getLoanRequestsList(stub, args)
		if err != nil {
			return nil, errors.New("Error getting loan requests list")
		}
		return res, nil
	}
	if function == "getLoanInvitationsQuantity" { //read a variable
		res, err := t.getLoanInvitationsQuantity(stub, args)
		if err != nil {
			return nil, errors.New("Error getting loan requests quantity")
		}
		return res, nil
	}
	if function == "getLoanInvitationsList" {
		res, err := t.getLoanInvitationsList(stub, args)
		if err != nil {
			return nil, errors.New("Error getting loan requests list")
		}
		return res, nil
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}
