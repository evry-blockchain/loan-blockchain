// SLSChainCode_test project SLSChainCode_test.go
package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, args []string, value string) {
	bytes, err := stub.MockQuery(name, args)
	if err != nil {
		fmt.Println("Query", name, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, function string, args []string) {
	_, err := stub.MockInvoke("1", "getBankId", args)
	if err != nil {
		fmt.Println("Invoke function", function, "with agrs", args, "failed", err)
		t.FailNow()
	}
}

func TestSLSChaincode_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=123 B=234
	checkInit(t, stub, []string{""})

	//checkState(t, stub, "A", "123")
	//checkState(t, stub, "B", "234")

	//addLoanRequest(stub, []string{"Statoil ASA", "6", "1M", "1M", "Statoil ASA project", "Statoil ASA project info", "Statoil ASA", "www.statoil.com", "John", "Smith", "10-01-2016", "Draft", "Oil industry"})
	/*checkQuery(t, stub, "A", "678")
	checkQuery(t, stub, "B", "567")
	checkState(t, stub, "A", "678")
	checkState(t, stub, "B", "567")	*/
}

/*func TestSLSChaincode_Query(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)
	checkQuery(t, stub, "countTableRowsg", []string{"Participants"}, "0")
}*/
