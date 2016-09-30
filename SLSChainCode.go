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
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	err := CreateParticipantTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating Participants table: " + err.Error())
	}

	err = CreateLoanTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating Loan table: " + err.Error())
	}
	err = CreateLoanSharesTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanShares table: " + err.Error())
	}
	err = CreateLoanRequestTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanRequests table: " + err.Error())
	}
	err = CreateLoanInvitationTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanInvitations table: " + err.Error())
	}
	err = CreateTransactionTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating Transactions table: " + err.Error())
	}
	err = CreateLoanReturnTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanReturn table: " + err.Error())
	}
	err = CreateLoanSaleTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanSale table: " + err.Error())
	}
	err = CreateLoanShareNegotiationTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating CreateLoanShareNegotiation table: " + err.Error())
	}
	err = createAccountTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating CreateLoanShareNegotiation table: " + err.Error())
	}

	t.populateInitialData(stub)

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	if function == "addParticipant" {
		return addParticipant(stub, args)
	}
	if function == "addLoan" {
		return addLoan(stub, args)
	}
	if function == "addLoanShare" {
		return addLoanShare(stub, args)
	}
	if function == "addLoanRequest" {
		return addLoanRequest(stub, args)
	}
	if function == "addLoanInvitation" {
		return addLoanInvitation(stub, args)
	}
	if function == "addTransaction" {
		return addTransaction(stub, args)
	}
	if function == "addLoanReturn" {
		return addLoanReturn(stub, args)
	}
	if function == "addLoanSale" {
		return addLoanSale(stub, args)
	}
	if function == "addLoanShareNegotiation" {
		return addLoanShareNegotiation(stub, args)
	}

	//Account
	if function == "addAccount" {
		return addAccount(stub, args)
	}
	if function == "updateAccountAmount" {
		return updateAccountAmount(stub, args)
	}

	// This function should not be invoked directly in the future!!!!!
	if function == "updateTableField" {
		return updateTableField(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	//Participants
	if function == "getParticipantsQuantity" {
		return getParticipantsQuantity(stub, args)
	}
	if function == "getParticipantsList" {
		return getParticipantsList(stub, args)
	}
	if function == "getParticipantsByType" {
		return getParticipantsByType(stub, args)
	}
	//Loans
	if function == "getLoansQuantity" {
		return getLoansQuantity(stub, args)
	}
	if function == "getLoansList" {
		return getLoansList(stub, args)
	}
	//LoanShares
	if function == "getLoanSharesQuantity" {
		return getLoanSharesQuantity(stub, args)
	}
	if function == "getLoanSharesList" {
		return getLoanSharesList(stub, args)
	}

	//LoanRequest
	if function == "getLoanRequestsQuantity" { //read a variable
		return getLoanRequestsQuantity(stub, args)
	}
	if function == "getLoanRequestsList" {
		return getLoanRequestsList(stub, args)
	}

	//LoanInvitation
	if function == "getLoanInvitationsQuantity" { //read a variable
		return getLoanInvitationsQuantity(stub, args)
	}
	if function == "getLoanInvitationsList" {
		return getLoanInvitationsList(stub, args)
	}

	//Transactions
	if function == "getTransactionsQuantity" { //read a variable
		return getTransactionsQuantity(stub, args)
	}
	if function == "getTransactionsList" {
		return getTransactionsList(stub, args)
	}

	//Loan Return
	if function == "getLoanReturnsQuantity" { //read a variable
		return getLoanReturnsQuantity(stub, args)
	}
	if function == "getLoanReturnsList" {
		return getLoanReturnsList(stub, args)
	}

	//Loan Sale
	if function == "getLoanSalesQuantity" { //read a variable
		return getLoanSalesQuantity(stub, args)
	}
	if function == "getLoanSalesList" {
		return getLoanSalesList(stub, args)
	}

	//Loan Share Negotiation
	if function == "getLoanShareNegotiationsQuantity" { //read a variable
		return getLoanShareNegotiationsQuantity(stub, args)
	}
	if function == "getLoanShareNegotiationsList" {
		return getLoanShareNegotiationsList(stub, args)
	}

	//Account
	if function == "getAccountsQuantity" { //read a variable
		return getAccountsQuantity(stub, args)
	}
	if function == "getAccountsList" {
		return getAccountsList(stub, args)
	}

	//========================================================================
	// This function should not be invoked directly in the future!!!!!!!!!!!!!!!!!!!!
	if function == "countTableRows" {
		return countTableRows(stub, args)
	}
	if function == "filterTableByValue" {
		return filterTableByValue(stub, args)
	}
	if function == "printCallerCertificate" {
		return printCallerCertificate(stub)
	}

	if function == "getCertAttribute" {
		return getCertAttribute(stub, args)
	}
	//========================================================================

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

func printCallerCertificate(stub shim.ChaincodeStubInterface) ([]byte, error) {
	// Verify the identity of the caller
	// Only an administrator can add Participant
	//###########################################

	certificate, err := stub.GetCallerCertificate()
	if err != nil {
		return nil, errors.New("Failed retrieving Certificate: " + err.Error())
	}
	fmt.Printf("\n\nCertificate: %v\n\n", string(certificate))

	callerMetadata, err := stub.GetCallerMetadata()
	if err != nil {
		return nil, errors.New("Failed retrieving Caller Metadata: " + err.Error())
	}
	fmt.Printf("Caller Metadata: %v\nCaller Metadata length:%v\n\n", string(callerMetadata), len(callerMetadata))

	payload, err := stub.GetPayload()
	if err != nil {
		return nil, errors.New("Failed retrieving Payload: " + err.Error())
	}
	fmt.Printf("Payload: %v\n\n", string(payload))

	binding, err := stub.GetBinding()
	if err != nil {
		return nil, errors.New("Failed retrieving Binding: " + err.Error())
	}
	fmt.Printf("Binding: %v\n\n", string(binding))

	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, errors.New("Failed retrieving Timestamp': " + err.Error())
	}
	fmt.Printf("Timestamp: %v\n\n", timestamp)

	//###########################################
	return []byte("Caller Metadata: " + string(callerMetadata)), nil
}

func getCertAttribute(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	attrName := args[0]
	attribute, err := stub.ReadCertAttribute(attrName)
	if err != nil {
		return nil, errors.New("Failed retrieving Certificate Attribute '" + attrName + "': " + err.Error())
	}
	fmt.Printf("Certificate Attribute '%v': %v\n\n", attrName, string(attribute))

	return []byte("Attribute '" + attrName + "': " + string(attribute)), nil
}

func (t *SimpleChaincode) populateInitialData(stub shim.ChaincodeStubInterface) error {

	//Participants
	_, _ = addParticipant(stub, []string{"Bank of Associates & Companies LTD", "Bank"})
	_, _ = addParticipant(stub, []string{"Connected Colaborators Bank", "Bank"})
	_, _ = addParticipant(stub, []string{"Bank of Paper, Wilson & Bluemine LTD", "Bank"})
	_, _ = addParticipant(stub, []string{"Bill Gates", "Borrower"})
	_, _ = addParticipant(stub, []string{"Peter Froystad", "Borrower"})
	_, _ = addParticipant(stub, []string{"John Smith", "Lowyer"})

	//Accounts
	_, _ = addAccount(stub, []string{"1", "10000"})
	_, _ = addAccount(stub, []string{"2", "50000"})
	_, _ = addAccount(stub, []string{"3", "30000"})

	//Loan Request
	// "BorrowerID", "LoanSharesAmount", "ProjectRevenue", "ProjectName", "ProjectInformation",
	//"Company", "Website", "ContactPersonName", "ContactPersonSurname", "RequestDate", "ArrangerBankID", "Status"
	_, _ = addLoanRequest(stub, []string{"1", "1000", "1M", "ProjectA", "ProjectA information", "CompanyA", "www.CompanyA.com", "Bill", "Gates", "10-01-2016", "1", "Pending"})
	_, _ = addLoanRequest(stub, []string{"1", "1000", "1M", "ProjectB", "ProjectB information", "CompanyB", "www.CompanyB.com", "Peter", "Froystad", "10-01-2016", "1", "Pending"})

	//Loan Invitation
	//"ArrangerBankID","BorrowerID","LoanRequestID","LoanTerm","Amount","InterestRate","Info","Status"
	_, _ = addLoanInvitation(stub, []string{"1", "1", "1", "2 years", "400", "3%", "Company A loan invitation info", "Pending"})
	_, _ = addLoanInvitation(stub, []string{"2", "3", "2", "3 years", "5000", "0.5%", "Company B loan invitation info", "Accepted"})

	//Loan Share Negotiation
	//"InvitationID","ParticipantBankID","Amount","NegotiationStatus"
	_, _ = addLoanShareNegotiation(stub, []string{"1", "2", "200", "Pending"})
	_, _ = addLoanShareNegotiation(stub, []string{"1", "3", "200", "Pending"})
	_, _ = addLoanShareNegotiation(stub, []string{"2", "1", "2000", "Pending"})
	_, _ = addLoanShareNegotiation(stub, []string{"2", "3", "3000", "Pending"})

	return nil
}
