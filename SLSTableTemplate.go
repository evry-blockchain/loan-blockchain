package main

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////// Uncomment below before using this table template ////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//Entity names
const <<X>>TableName = ""

//Column names
const <<Xpref>>_<<Column1>>ColName = "<<Column1>>"
const <<Xpref>>_<<Column2>>ColName = "<<Column2>>"
const <<Xpref>>_<<Column3>>ColName = "<<Column3>>"

//Column quantity
const <<X>>TableColsQty = <<Write columns quantity here>>

// ============================================================================================================================
//
// ============================================================================================================================

func Create<<X>>Table(stub shim.ChaincodeStubInterface) error {
	<<Xpref>>_ColumnNames := []string{<<Xpref>>_<<Column1>>ColName, <<Xpref>>_<<Column2>>ColName, <<Xpref>>_<<Column3>>ColName}
	return createTable(stub, <<X>>TableName, <<Xpref>>_ColumnNames)
}

func add<<X>>(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	attrName := "role"
	attrValue := "assigner"
	checkPermissionsAssigner, errA := checkAttribute(stub, attrName, attrValue)
	if !checkPermissionsAssigner {
		return nil, errors.New("Error checking permission in add<<X>>: " + errA.Error())
	}

	if len(args) == <<X>>TableColsQty {
		return nil, addRow(stub, <<X>>TableName, args, true)
	}
	if len(args) == <<X>>TableColsQty-1 {
		return nil, addRow(stub, <<X>>TableName, args, false)
	}

	return nil, errors.New("Incorrect number of arguments. " +
		"Provided " + strconv.Itoa(len(args)) + "Expecting " + strconv.Itoa(<<X>>TableColsQty-1) +
		" or " + strconv.Itoa(<<X>>TableColsQty))
}

func get<<X>>Quantity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return countTableRows(stub, []string{<<X>>TableName})
}

func get<<X>>List(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return filterTableByValue(stub, []string{<<X>>TableName})
}


func get<<X>>ByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	keyValue := args[0]
	return filterTableByKey(stub, <<X>>TableName, keyValue)
}

func get<<X>>MaxKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxKey, err := getTableMaxKey(stub, <<X>>TableName)
	if err != nil {
		return nil, errors.New("Error in get<<X>>MaxKey func: " + err.Error())
	}
	return maxKey, nil
}

func update<<X>>(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != <<X>>TableColsQty {
		return nil, errors.New("Incorrect number of arguments in update<<X>> func. Expecting " + strconv.Itoa(<<X>>TableColsQty))
	}

	tbl, err := stub.GetTable(<<X>>TableName)
	if err != nil {
		return nil, errors.New("An error occured while running update<<X>>: " + err.Error())
	}

	for i, cd := range tbl.ColumnDefinitions {
		_, err := updateTableField(stub, []string{<<X>>TableName, args[0], cd.Name, args[i]})		//args[0] is hardcoded as row id
		if err != nil {
			return nil, errors.New("Failed updating field '" + cd.Name + "' in update<<X>> func: " + err.Error())
		}
	}

	return nil, nil
}

*/
