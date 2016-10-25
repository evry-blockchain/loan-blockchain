package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
//
// ============================================================================================================================
func getProjectsList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	bankid, err := getBankId(stub, []string{})
	if err != nil {
		return nil, errors.New("Error getting bankid in getProjectsList func: " + err.Error())
	}

	/*
		maxKey, err := getTableMaxKey(stub, LoanNegotiationsTableName)
		if err != nil {
			return nil, errors.New("Error getting maxKey in getProjectsList func: " + err.Error())
		}
	*/
	return filterTableByValue(stub, []string{LoanRequestsTableName, LR_ArrangerBankIDColName, string(bankid)})
}
