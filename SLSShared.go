package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// This function should update single column of any table
// If column or table does not exists it returns error
// If filter takes quantity of rows not equal to one (zero as well) it returns error
func updateTableField(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	tableName, keyValue, columnName, columnNewValue := args[0], args[1], args[2], args[3]

	tbl, err := stub.GetTable(tableName)
	if err != nil {
		return nil, errors.New("An error occured while running updateTableField: " + err.Error())
	}

	//Delete in the future
	fmt.Println("Table '" + tbl.Name + "' exists. Trying to update single value in '" + columnName + "' column.")

	var cols []shim.Column
	col := shim.Column{Value: &shim.Column_String_{String_: keyValue}}
	cols = append(cols, col)

	row, errR := stub.GetRow(tableName, cols)
	if errR != nil {
		return nil, errors.New("An error occured while running updateTableField: " + errR.Error())
	}
	if row.GetColumns() == nil {
		return nil, errors.New("Key value not found")
	}

	var f bool
	var columnOldValue string

	// This Println is temporary for logging purposes
	fmt.Println("====== Before update =========================")
	for i, c := range row.GetColumns() {
		// This Printf is temporary for logging purposes
		fmt.Printf("Column number: '%v', Column name: '%v', Column value: '%v'\n", i, tbl.ColumnDefinitions[i].Name, c.GetString_())
		if tbl.ColumnDefinitions[i].Name == columnName {
			columnOldValue = c.GetString_()
			row.Columns[i] = &shim.Column{Value: &shim.Column_String_{String_: columnNewValue}}
			f = true
		}
	}

	// This Println is temporary for logging purposes
	fmt.Println("====== After update =========================")
	// This Printf is temporary for logging purposes
	for i, c := range row.GetColumns() {
		fmt.Printf("Column number: '%v', Column name: '%v', Column value: '%v'\n", i, tbl.ColumnDefinitions[i].Name, c.GetString_())
	}
	// This Println is temporary for logging purposes
	fmt.Println("=============================================")

	if !f {
		return nil, errors.New("Column '" + columnName + "' is missing")
	}

	ok, errreplace := stub.ReplaceRow(tableName, row)
	if errreplace != nil {
		return nil, errors.New("An error occured while running updateTableField: " + errreplace.Error())
	}
	//This check might be redundant.
	if !ok {
		return nil, errors.New("A row does not exist the given key")
	}

	fmt.Printf("Column '%v' of the row of key value '%v' in the table '%v' has been successfuly updated from value '%v' to value '%v'\n", columnName, keyValue, tableName, columnOldValue, columnNewValue)

	return nil, nil
}

func countTableRows(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var numberOfArgs int = 1
	if len(args) != numberOfArgs {
		return nil, errors.New("Incorrect number of arguments. Expecting: " + strconv.Itoa(numberOfArgs))
	}

	tableName := args[0]

	q, err := countTableRowsInt(stub, tableName)
	if err != nil {
		return nil, errors.New("Failed to get rows quantity for table '" + tableName + "': " + err.Error())
	}

	fmt.Printf("Quantity of rows for table %v is %v\n", tableName, q)
	return []byte(strconv.Itoa(q)), nil
}

func countTableRowsInt(stub shim.ChaincodeStubInterface, tableName string) (int, error) {
	// The function hangs for about 10 seconds if table Name does not exist
	// consider a fix !!!!!!!!!!!!!!

	// Use emty columns slice to get all rows for count
	var cols []shim.Column

	rowChan, err := stub.GetRows(tableName, cols)
	if err != nil {
		return 0, err
	}

	_, ok := <-rowChan

	var q int
	for ok {
		//		fmt.Printf("ok: %v\n", ok)
		//		fmt.Printf("Rows to string: %v\n", row2)
		q++
		_, ok = <-rowChan
	}

	return q, nil
}

func filterTableByValue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var tableName, filterColumn, filterValue string
	var isFiltered bool
	// 1 or 3 arguments should be provided:
	// 1 when filter is not needed: table name
	// 3 when filter is need: table name, filter column, filter value

	switch l := len(args); l {
	case 1:
		tableName = args[0]
		isFiltered = false
	case 3:
		tableName, filterColumn, filterValue = args[0], args[1], args[2]
		isFiltered = true
	default:
		return nil, errors.New("Incorrect number of arguments. Expecting: 1 or 3")
	}

	tbl, err := stub.GetTable(tableName)
	if err != nil {
		return nil, errors.New("An error occured while running filterTableByValue: " + err.Error())
	}

	var cols []shim.Column
	var rows []shim.Row

	rowChan, _ := stub.GetRows(tableName, cols)
	row, ok := <-rowChan

	if isFiltered {
		var columnNumber int
		for i, cd := range tbl.ColumnDefinitions {
			if cd.Name == filterColumn {
				columnNumber = i
			}
		}

		var i int
		for ok {
			if row.Columns[columnNumber].GetString_() == filterValue {
				rows = append(rows, row)
				i++
			}
			row, ok = <-rowChan
		}
	} else {
		for ok {
			rows = append(rows, row)
			row, ok = <-rowChan
		}
	}

	//fmt.Printf("Filter table %v column %v by value %v\n Result:\nRows quantity: %v\n Rows content: %v\n", tableName, filterColumn, filterValue, i, rows)

	return recordsetToJson(stub, tbl, rows)
}

func recordsetToJson(stub shim.ChaincodeStubInterface, tbl *shim.Table, rows []shim.Row) ([]byte, error) {

	var ColumnNames []string
	for _, cd := range tbl.ColumnDefinitions {
		ColumnNames = append(ColumnNames, cd.Name)
	}

	var s string
	s = "["

	for _, r := range rows {
		s += "{"
		for m, c := range r.Columns {

			columnName := tbl.ColumnDefinitions[m].Name
			columnValue := c.GetString_()
			s += "\"" + columnName + "\":\"" + columnValue + "\","
		}
		s = s[0:len(s)-1] + "},"
	}

	s = s[0:len(s)-1] + "]"

	return []byte(s), nil
}

func createTable(stub shim.ChaincodeStubInterface, tableName string, columns []string) error {
	//not to forget delete table is it already exists
	stub.DeleteTable(tableName)

	var colDefs []*shim.ColumnDefinition

	for _, colName := range columns {
		colDefs = append(colDefs, &shim.ColumnDefinition{Name: colName, Type: shim.ColumnDefinition_STRING, Key: false})
	}
	colDefs[0].Key = true

	err := stub.CreateTable(tableName, colDefs)
	if err != nil {
		return errors.New("Failed to add table '" + tableName + "' to state: " + err.Error())
	}
	fmt.Println("Table '" + tableName + "' created")
	return nil
}

func addRow(stub shim.ChaincodeStubInterface, tableName string, args []string) error {

	q, err := countTableRowsInt(stub, tableName)
	if err != nil {
		return errors.New("Failed to add row to '" + tableName + "' table: " + err.Error())
	}
	q++
	qstr := strconv.Itoa(q)

	tbl, _ := stub.GetTable(tableName)
	colDefs := tbl.ColumnDefinitions
	colsQty := len(colDefs)

	//Add Account to ledger table
	var cols []*shim.Column
	cols = append(cols, &shim.Column{Value: &shim.Column_String_{String_: qstr}})
	for i := 1; i <= colsQty-1; i++ {
		cols = append(cols, &shim.Column{Value: &shim.Column_String_{String_: args[i-1]}})
	}

	var ok bool
	ok, err = stub.InsertRow(tableName, shim.Row{Columns: cols})
	if err != nil {
		return errors.New("Failed to add row to '" + tableName + "' table: " + err.Error())
	}
	if !ok {
		return errors.New("Row with key '" + qstr + "' is already assigned in table '" + tableName + "'")
	}

	s := "The row has been added to table '" + tableName + "' in ledger: \n"
	for i, cd := range colDefs {
		s += cd.Name + " : " + cols[i].GetString_() + "\n"
	}

	fmt.Println(s)
	return nil

}
