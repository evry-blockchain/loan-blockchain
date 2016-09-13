package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Entity names
const LoanRequestsTableName = "LoanRequests"
const LoanRequestsQuantityKey = "LoanRequestsQuantity"

//Column names
const LR_LoanRequestIDColName = "LoanRequestID"
const LR_BorrowerIDColName = "BorrowerID"
const LR_LoanSharesAmountColName = "LoanSharesAmount"
const LR_ProjectRevenueColName = "ProjectRevenue"
const LR_ProjectNameColName = "ProjectName"
const LR_ProjectInformationColName = "ProjectInformation"
const LR_CompanyColName = "Company"
const LR_WebsiteColName = "Website"
const LR_ContactPersonNameColName = "ContactPersonName"
const LR_ContactPersonSurnameColName = "ContactPersonSurname"
const LR_RequestDateColName = "RequestDate"
const LR_ArrangerBankIDColName = "ArrangerBankID"
const LR_StatusColName = "Status"

const LoanRequestColumnsQuantity = 13

//struct for json (consider avoid it)
type LoanRequest struct {
	LoanRequestID        string
	BorrowerID           string
	LoanSharesAmount     string
	ProjectRevenue       string
	ProjectNameCol       string
	ProjectInformation   string
	Company              string
	Website              string
	ContactPersonName    string
	ContactPersonSurname string
	RequestDate          string
	ArrangerBankID       string
	Status               string
}

var LoanRequestsQuantity int

// ============================================================================================================================
//
// ============================================================================================================================

func CreateLoanRequestTable(stub *shim.ChaincodeStub) error {
	stub.DeleteTable(LoanRequestsTableName)

	cdfs := []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: LR_LoanRequestIDColName, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: LR_BorrowerIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_LoanSharesAmountColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_ProjectRevenueColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_ProjectNameColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_ProjectInformationColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_CompanyColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_WebsiteColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_ContactPersonNameColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_ContactPersonSurnameColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_RequestDateColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_ArrangerBankIDColName, Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: LR_StatusColName, Type: shim.ColumnDefinition_STRING, Key: false},
	}

	err := stub.CreateTable(LoanRequestsTableName, cdfs)
	if err != nil {
		return errors.New("Failed to add Loan Requests table to state: " + err.Error())
	}

	err2 := stub.PutState(LoanRequestsQuantityKey, []byte(strconv.Itoa(0)))
	if err2 != nil {
		return errors.New("Failed to add Loan Requests quantity to state: " + err2.Error())
	}

	LoanRequestsQuantity = 0

	return nil
}

func (t *SimpleChaincode) addLoanRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != LoanRequestColumnsQuantity-1 {
		return nil, errors.New("Incorrect number of arguments. Expecting " + strconv.Itoa(LoanRequestColumnsQuantity-1))
	}

	qBytes, _ := stub.GetState(LoanRequestsQuantityKey)
	if qBytes == nil {
		return nil, errors.New("Loan Requests quantity entity not found")
	}

	qstr := string(qBytes)
	q, errconv := strconv.Atoi(qstr)
	if errconv != nil {
		return nil, errors.New("Error converting key string to int")
	}
	q++
	qstr = strconv.Itoa(q)

	//Add Loan to ledger table
	var cols []*shim.Column
	cols = append(cols, &shim.Column{Value: &shim.Column_String_{String_: qstr}})
	for i := 1; i <= LoanRequestColumnsQuantity-1; i++ {
		cols = append(cols, &shim.Column{Value: &shim.Column_String_{String_: args[i-1]}})
	}
	ok, err := stub.InsertRow(LoanRequestsTableName, shim.Row{Columns: cols})
	if !ok && err == nil {
		return nil, errors.New("Loan with LoanID " + qstr + " is already assigned")
	}

	//Update total Loan Requests quantity
	err2 := stub.PutState(LoanRequestsQuantityKey, []byte(qstr))

	LoanRequestsQuantity = q

	if err != nil || err2 != nil {
		return nil, errors.New("Failed to add Loan")
	}

	s := "The row has been added to Loan Requests table in ledger: \n" +
		LR_LoanRequestIDColName + " : " + qstr + "\n" +
		LR_BorrowerIDColName + " : " + args[0] + "\n" +
		LR_LoanSharesAmountColName + " : " + args[1] + "\n" +
		LR_ProjectRevenueColName + " : " + args[2] + "\n" +
		LR_ProjectNameColName + " : " + args[3] + "\n" +
		LR_ProjectInformationColName + " : " + args[4] + "\n" +
		LR_CompanyColName + " : " + args[5] + "\n" +
		LR_WebsiteColName + " : " + args[6] + "\n" +
		LR_ContactPersonNameColName + " : " + args[7] + "\n" +
		LR_ContactPersonSurnameColName + " : " + args[8] + "\n" +
		LR_RequestDateColName + " : " + args[9] + "\n" +
		LR_ArrangerBankIDColName + " : " + args[10] + "\n" +
		LR_StatusColName + " : " + args[11]

	fmt.Println(s)
	return nil, nil
}

func (t *SimpleChaincode) getLoanRequestsQuantity(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return []byte(strconv.Itoa(LoanRequestsQuantity)), nil
}

func (t *SimpleChaincode) getLoanRequestsList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var p LoanRequest
	var ps []LoanRequest

	for i := 1; i <= LoanRequestsQuantity; i++ {
		var cols []shim.Column
		col := shim.Column{Value: &shim.Column_Int32{Int32: int32(i)}}
		cols = append(cols, col)

		row, _ := stub.GetRow(LoanRequestsTableName, cols)

		p.LoanRequestID = row.Columns[0].GetString_()
		p.BorrowerID = row.Columns[1].GetString_()
		p.LoanSharesAmount = row.Columns[2].GetString_()
		p.ProjectRevenue = row.Columns[3].GetString_()
		p.ProjectNameCol = row.Columns[4].GetString_()
		p.ProjectInformation = row.Columns[5].GetString_()
		p.Company = row.Columns[6].GetString_()
		p.Website = row.Columns[7].GetString_()
		p.ContactPersonName = row.Columns[8].GetString_()
		p.ContactPersonSurname = row.Columns[9].GetString_()
		p.RequestDate = row.Columns[10].GetString_()
		p.ArrangerBankID = row.Columns[11].GetString_()
		p.Status = row.Columns[12].GetString_()

		ps = append(ps, p)
	}

	res, _ := json.Marshal(ps)

	return res, nil
}
