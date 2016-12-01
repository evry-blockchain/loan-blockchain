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

const isAuthenticationEnabled = false

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
	err = CreateLoanRequestTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanRequests table: " + err.Error())
	}
	err = CreateLoanNegotiationTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanNegotiations table: " + err.Error())
	}
	err = CreateLoanTermTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanTerms table: " + err.Error())
	}
	err = CreateLoanTermProposalTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanTermProposals table: " + err.Error())
	}
	err = CreateLoanTermVoteTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanTermVotes table: " + err.Error())
	}
	err = CreateLoanTermCommentTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating LoanTermComments table: " + err.Error())
	}
	err = CreateUserTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating Users table: " + err.Error())
	}

	populateInitialData(stub, args)

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}

	//========================================================================
	//Participant
	if function == "addParticipant" {
		return addParticipant(stub, args)
	}

	//========================================================================
	//LoanRequest
	if function == "addLoanRequest" {
		return addLoanRequest(stub, args)
	}
	if function == "updateLoanRequest" {
		return updateLoanRequest(stub, args)
	}

	//========================================================================
	//Loan Negotiation
	if function == "addLoanNegotiation" {
		return addLoanNegotiation(stub, args)
	}
	if function == "updateLoanNegotiation" {
		return updateLoanNegotiation(stub, args)
	}
	if function == "updateLoanNegotiationStatus" {
		return updateLoanNegotiationStatus(stub, args)
	}
	if function == "updateParticipantBankComment" {
		return updateParticipantBankComment(stub, args)
	}

	//========================================================================
	//Loan Term
	if function == "addLoanTerm" {
		return addLoanTerm(stub, args)
	}
	if function == "updateLoanTerm" {
		return updateLoanTerm(stub, args)
	}

	//========================================================================
	//Loan Term Proposal
	if function == "addLoanTermProposal" {
		return addLoanTermProposal(stub, args)
	}
	if function == "updateLoanTermProposal" {
		return updateLoanTermProposal(stub, args)
	}

	//========================================================================
	//Loan Term Proposal
	if function == "addLoanTermProposal" {
		return addLoanTermProposal(stub, args)
	}
	if function == "updateLoanTermProposal" {
		return updateLoanTermProposal(stub, args)
	}

	//========================================================================
	//Loan Term Vote
	if function == "addLoanTermVote" {
		return addLoanTermVote(stub, args)
	}
	if function == "updateLoanTermVote" {
		return updateLoanTermVote(stub, args)
	}

	//========================================================================
	//Loan Term Comment
	if function == "addLoanTermComment" {
		return addLoanTermComment(stub, args)
	}
	if function == "updateLoanTermComment" {
		return updateLoanTermComment(stub, args)
	}

	//========================================================================
	//User
	if function == "addUser" {
		return addUser(stub, args)
	}
	if function == "updateUser" {
		return updateUser(stub, args)
	}

	//========================================================================
	// Specific functions
	if function == "updateTableField" {
		return updateTableField(stub, args)
	}
	if function == "deleteRow" {
		return deleteRow(stub, args)
	}
	if function == "deleteRowsByColumnValue" {
		return deleteRowsByColumnValue(stub, args)
	}
	if function == "populateInitialData" {
		return populateInitialData(stub, args)
	}
	//========================================================================

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	//========================================================================
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
	if function == "getParticipantsByKey" {
		return getParticipantsByKey(stub, args)
	}
	if function == "getParticipantsMaxKey" {
		return getParticipantsMaxKey(stub, args)
	}

	//========================================================================
	//LoanRequest
	if function == "getLoanRequestsQuantity" {
		return getLoanRequestsQuantity(stub, args)
	}
	if function == "getLoanRequestsList" {
		return getLoanRequestsList(stub, args)
	}
	if function == "getLoanRequestByKey" {
		return getLoanRequestByKey(stub, args)
	}
	if function == "getLoanRequestsMaxKey" {
		return getLoanRequestsMaxKey(stub, args)
	}

	//========================================================================
	//Loan Negotiation
	if function == "getLoanNegotiationsQuantity" {
		return getLoanNegotiationsQuantity(stub, args)
	}
	if function == "getLoanNegotiationsList" {
		return getLoanNegotiationsList(stub, args)
	}
	if function == "getLoanNegotiationByKey" {
		return getLoanNegotiationByKey(stub, args)
	}
	if function == "getLoanNegotiationsMaxKey" {
		return getLoanNegotiationsMaxKey(stub, args)
	}

	//========================================================================
	//Loan Term
	if function == "getLoanTermQuantity" {
		return getLoanTermQuantity(stub, args)
	}
	if function == "getLoanTermList" {
		return getLoanTermList(stub, args)
	}
	if function == "getLoanTermByKey" {
		return getLoanTermByKey(stub, args)
	}
	if function == "getLoanTermMaxKey" {
		return getLoanTermMaxKey(stub, args)
	}

	//========================================================================
	//Loan Term Proposal
	if function == "getLoanTermProposalQuantity" {
		return getLoanTermProposalQuantity(stub, args)
	}
	if function == "getLoanTermProposalList" {
		return getLoanTermProposalList(stub, args)
	}
	if function == "getLoanTermProposalByKey" {
		return getLoanTermProposalByKey(stub, args)
	}
	if function == "getLoanTermProposalMaxKey" {
		return getLoanTermProposalMaxKey(stub, args)
	}

	//========================================================================
	//Loan Term Vote
	if function == "getLoanTermVoteQuantity" {
		return getLoanTermVoteQuantity(stub, args)
	}
	if function == "getLoanTermVoteList" {
		return getLoanTermVoteList(stub, args)
	}
	if function == "getLoanTermVoteByKey" {
		return getLoanTermVoteByKey(stub, args)
	}
	if function == "getLoanTermVoteMaxKey" {
		return getLoanTermVoteMaxKey(stub, args)
	}

	//========================================================================
	//Loan Term Comment
	if function == "getLoanTermCommentQuantity" {
		return getLoanTermCommentQuantity(stub, args)
	}
	if function == "getLoanTermCommentList" {
		return getLoanTermCommentList(stub, args)
	}
	if function == "getLoanTermCommentByKey" {
		return getLoanTermCommentByKey(stub, args)
	}
	if function == "getLoanTermCommentMaxKey" {
		return getLoanTermCommentMaxKey(stub, args)
	}

	//========================================================================
	//User
	if function == "getUserQuantity" {
		return getUserQuantity(stub, args)
	}
	if function == "getUserList" {
		return getUserList(stub, args)
	}
	if function == "getUserByKey" {
		return getUserByKey(stub, args)
	}
	if function == "getUserMaxKey" {
		return getUserMaxKey(stub, args)
	}

	//========================================================================
	// Special functions
	if function == "countTableRows" {
		return countTableRows(stub, args)
	}
	if function == "filterTableByValue" {
		return filterTableByValue(stub, args)
	}
	/*if function == "printCallerCertificate" {
		return printCallerCertificate(stub)
	}*/
	if function == "getCertAttribute" {
		return getCertAttribute(stub, args)
	}
	if function == "getBankId" {
		return getBankId(stub, args)
	}
	if function == "getProjectsList" {
		return getProjectsList(stub, args)
	}
	//========================================================================

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

func getCertAttribute(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments in getCertAttribute func. Expecting 1")
	}

	attrName := args[0]
	attribute, err := stub.ReadCertAttribute(attrName)
	if err != nil {
		return nil, errors.New("Failed retrieving Certificate Attribute '" + attrName + "' in getCertAttribute func: " + err.Error())
	}

	return []byte("Attribute '" + attrName + "': " + string(attribute)), nil
}

func getBankId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments in getBankId func. Expecting 0")
	}

	attrName := "bankid"
	attribute, err := stub.ReadCertAttribute(attrName)
	if err != nil {
		return nil, errors.New("Failed retrieving Certificate Attribute '" + attrName + "' in getBankId func: " + err.Error())
	}

	return []byte(string(attribute)), nil
}

func checkAttribute(stub shim.ChaincodeStubInterface, attrName, attrValue string) (bool, error) {
	if !isAuthenticationEnabled {
		return true, nil
	}
	// Why stub.VerifyAttribute is not used here?????? Consider using it.
	attribute, err := stub.ReadCertAttribute(attrName)
	if err != nil {
		return false, errors.New("Error checking role: " + err.Error())
	}
	if string(attribute) != attrValue {
		return false, errors.New("Current user attribute '" + attrName + "' value is '" + string(attribute) + "' but not '" + attrValue + "'")
	}
	return true, nil
}

func populateInitialData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//Participants
	_, _ = deleteRowsByColumnValue(stub, []string{ParticipantsTableName})
	//Adding banks with keys
	_, _ = addParticipant(stub, []string{"6", "SpareBank 1 SR-BANK", "Bank"})
	_, _ = addParticipant(stub, []string{"7", "DNB ASA", "Bank"})
	_, _ = addParticipant(stub, []string{"8", "Nationwide Building Society", "Bank"})
	_, _ = addParticipant(stub, []string{"9", "JPMorgan Chase & Co", "Bank"})
	_, _ = addParticipant(stub, []string{"10", "Barclays", "Bank"})
	_, _ = addParticipant(stub, []string{"11", "Mizuho Bank, Ltd.", "Bank"})

	//Accounts
	_, _ = deleteRowsByColumnValue(stub, []string{AccountsTableName})
	_, _ = addAccount(stub, []string{"1", "10000"})
	_, _ = addAccount(stub, []string{"2", "50000"})
	_, _ = addAccount(stub, []string{"3", "30000"})

	//Loan Request
	// "BorrowerID", "ArrangerBankID", "LoanSharesAmount", "ProjectRevenue", "ProjectName", "ProjectInformation",
	//"Company", "Website", "ContactPersonName", "ContactPersonSurname", "RequestDate",
	//"Status", "MarketAndIndustry"
	_, _ = deleteRowsByColumnValue(stub, []string{LoanRequestsTableName})
	_, _ = addLoanRequest(stub, []string{"Statoil ASA", "6", "1M", "1M", "Statoil ASA project",
		"Statoil ASA project info", "Statoil ASA", "www.statoil.com",
		"John", "Smith", "10-01-2016", "Draft", "Oil industry",
		"some LoanTerm", "some Assets", "some Convenants", "some InterestRate"})
	_, _ = addLoanRequest(stub, []string{"BP Global", "7", "1M", "1M", "BP Global project",
		"BP Global project info", "BP Global", "www.bp.com", "Peter",
		"Froystad", "10-01-2016", "Draft", "Oil industry",
		"some LoanTerm", "some Assets", "some Convenants", "some InterestRate"})

	//Loan Share Negotiation
	//"InvitationID","ParticipantBankID","Amount","NegotiationStatus", "ParticipantBankComment", "Date"
	_, _ = deleteRowsByColumnValue(stub, []string{LoanNegotiationsTableName})
	_, _ = addLoanNegotiation(stub, []string{"1", "6", "200 M USD", "INVITED", "Comment of SpareBank", "11-01-2016"})
	_, _ = addLoanNegotiation(stub, []string{"1", "9", "100 M USD", "INVITED", "Comment of JPMorgan", "12-01-2016"})
	_, _ = addLoanNegotiation(stub, []string{"1", "10", "100 M USD", "INVITED", "Comment of Barclays", "12-01-2016"})
	_, _ = addLoanNegotiation(stub, []string{"2", "7", "250 M USD", "INVITED", "Comment of Nationwide Building Society", "21-01-2016"})
	_, _ = addLoanNegotiation(stub, []string{"2", "9", "200 M USD", "INVITED", "Comment of JPMorgan", "22-01-2016"})
	_, _ = addLoanNegotiation(stub, []string{"2", "11", "300 M USD", "INVITED", "Comment of Mizuho Bank, Ltd.", "22-01-2016"})

	return nil, nil
}
