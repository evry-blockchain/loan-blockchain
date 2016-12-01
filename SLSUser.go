package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const UserTableName = "Users"

//Column names
const U_UserIDColName = "UserID"
const U_ParticipantIDColName = "ParticipantID"
const U_UserNameColName = "UserName"

//Column quantity
const UserTableColsQty = 3

// ============================================================================================================================
//
// ============================================================================================================================

func CreateUserTable(stub shim.ChaincodeStubInterface) error {
	U_ColumnNames := []string{U_UserIDColName, U_ParticipantIDColName, U_UserNameColName}
	return createTable(stub, UserTableName, U_ColumnNames)
}

func addUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission in addUser: " + errA.Error())
	}

	if len(args) == UserTableColsQty {
		return nil, addRow(stub, UserTableName, args, true)
	}
	if len(args) == UserTableColsQty-1 {
		return nil, addRow(stub, UserTableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + "Expecting " + strconv.Itoa(UserTableColsQty-1) +
		" or " + strconv.Itoa(UserTableColsQty))
}

func getUserQuantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{UserTableName})
}

func getUserList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{UserTableName})
}

func getUserByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, UserTableName, keyValue)
}

func getUserMaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, UserTableName)
	if err != nil {
		return nil, errors.New("Error in getUserMaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func updateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != UserTableColsQty {
		return nil, errors.New("Incorrect number of arguments in updateUser func. Expecting " + strconv.Itoa(UserTableColsQty))
	}

	tbl, err := stub.GetTable(UserTableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateUser: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{UserTableName, args[0], cd.Name, args[i]}) //args[0] is hardcoded as row id
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in updateUser func: " + err.Error())
		}
	}

	return nil, nil
}
