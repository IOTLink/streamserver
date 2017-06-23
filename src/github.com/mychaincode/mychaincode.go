package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init ...
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Init ###########")
	_, args := stub.GetFunctionAndParameters()
	var A string
	var Aval int64 // Asset holdings
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	A = args[0]
	// Initialize the chaincode
	Aval, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d\n", Aval)
	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.FormatInt(Aval,10)))
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			fmt.Printf("Transient data in 'init' : %s\n", transientData)
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

// Query ...
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var A string
	A = args[0]
	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

// Invoke ...
// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "move" {
		if err := stub.SetEvent("testEvent", []byte("Test Payload")); err != nil {
			return shim.Error("Unable to set CC event: testEvent. Aborting transaction ...")
		}
		return t.move(stub, args)
	}

	if args[0] == "init" {
		return t.init(stub, args)
	}

	if args[0] == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}

	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}

func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var A, B string    // Entities
	var Aval, Bval int64 // Asset holdings
	var X int64          // Transaction value
	var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	A = args[1]
	B = args[2]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		err = stub.PutState(A, []byte("0"))
		if err != nil {
			return shim.Error("Entity not found")
			//return shim.Error(err.Error())
		}
		Avalbytes = []byte("0")
	}
	Aval, _ = strconv.ParseInt(string(Avalbytes), 10, 64)

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		err = stub.PutState(B, []byte("0"))
		if err != nil {
			return shim.Error("Entity not found")
			//return shim.Error(err.Error())
		}
		Bvalbytes = []byte("0")
	}
	Bval, _ = strconv.ParseInt(string(Bvalbytes), 10, 64)

	// Perform the execution
	X, err = strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.FormatInt(Aval,10)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.FormatInt(Bval,10)))
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			fmt.Printf("Transient data in 'move' : %s\n", transientData)
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[1]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) init(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var A string
	var Aval int64 // Asset holdings
	var err error

	A = args[1]
	// Initialize the chaincode
	Aval, err = strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d\n", Aval)
	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.FormatInt(Aval,10)))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[1]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

