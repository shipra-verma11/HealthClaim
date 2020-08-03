package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	// "github.com/hyperledger/fabric-chaincode-go/shim"
	// pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	policyCounterNo   = "PolicyCounterNo"
	customerCounterNo = "CustomerCounterNo"
)

type insuranceClaim struct {
}

type policyAsset struct {
	PolicyNumber     string
	CustomerID       int
	ClaimLastUpdated time.Time
	ClaimStatus      string
}

type customerAsset struct {
	ID   int
	Name string
	Age  string
	Sex  string
	Pan  string
}

type counterNo struct {
	Counter int
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(insuranceClaim))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}

}

// ============================================================================================================================
// Init
// ============================================================================================================================
func (t *insuranceClaim) Init(APIstub shim.ChaincodeStubInterface) pb.Response {

	customerCounterAsset, _ := APIstub.GetState(customerCounterNo)
	if customerCounterAsset == nil {
		customerCounterByte, _ := json.Marshal(counterNo{Counter: 0})
		err := APIstub.PutState(customerCounterNo, customerCounterByte)
		if err != nil {
			return shim.Error("Couldn't able to update customer counter")
		}
	}

	policyCounterAsset, _ := APIstub.GetState(policyCounterNo)
	if policyCounterAsset == nil {
		policyCounterByte, _ := json.Marshal(counterNo{Counter: 56})
		err := APIstub.PutState(policyCounterNo, policyCounterByte)
		if err != nil {
			return shim.Error("Couldn't able to update policy counter")
		}
	}

	return shim.Success([]byte("Insurace Claim CC initialized successfully"))
}

// ============================================================================================================================
// Invoke
// ============================================================================================================================
func (t *insuranceClaim) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {

	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("Invoke function ==> :" + function)
	if function == "registerCustomer" {
		return t.registerCustomer(APIstub, args)
	} else if function == "issuePolicy" {
		return t.issuePolicy(APIstub, args)
	} else if function == "queryCustomerRecord" {
		return t.queryCustomerRecord(APIstub, args)
	} else if function == "claim" {
		return t.claim(APIstub, args)
	} else if function == "queryAllAsset" {
		return t.queryAllAsset(APIstub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function:" + function)
}

// register the customer
// Param: Name, Age, Sex, Pan
func (t *insuranceClaim) registerCustomer(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************************** registerCustomer")
	fmt.Println("no.of arguments:", args)

	// Check length of the arguments
	validateArgs, errString := validateArgs(args, 4)
	if !validateArgs {
		return shim.Error(errString)
	}

	// check whether customer exists
	costumerD, err := APIstub.GetState(args[3])
	if err != nil {
		return shim.Error("Function registerCustomer failed to execute, Err:" + err.Error())
	}
	if costumerD == nil {
		// get customer id
		getCustomerID := getCounter(APIstub, customerCounterNo)
		// upload customer data
		registerUserDate := customerAsset{
			ID:   getCustomerID + int(1),
			Name: args[0],
			Age:  args[1],
			Sex:  args[2],
			Pan:  args[3],
		}
		assetAsBytes, errMarshal := json.Marshal(registerUserDate)
		if errMarshal != nil {
			return shim.Error(fmt.Sprintf("Not able to marshal for user record: %s", registerUserDate.Pan))
		}

		errPut := APIstub.PutState(registerUserDate.Pan, assetAsBytes)
		fmt.Println(errPut)
		if errPut != nil {
			return shim.Error(fmt.Sprintf("Failed to register User: %s", registerUserDate.Pan))
		}

		customerCounterAsset := counterNo{
			Counter: registerUserDate.ID,
		}

		counterAsBytes, _ := json.Marshal(customerCounterAsset)
		errCounter := APIstub.PutState(customerCounterNo, counterAsBytes)
		if errCounter != nil {
			return shim.Error("Failed to update the customer counter")
		}

		return shim.Success([]byte(fmt.Sprintf("Customer ID:`%d` successfully registered, for PAN `%s`", registerUserDate.ID, registerUserDate.Pan)))
	}

	return shim.Error(fmt.Sprintf("Customer record already exists: %s", args))
}

// queryCustomerRecord
// Param: Pan
func (t *insuranceClaim) queryCustomerRecord(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************************** queryCustomerRecord")
	fmt.Println("no.of arguments:", args)

	// Check length of the arguments
	validateArgs, errString := validateArgs(args, 1)
	if !validateArgs {
		return shim.Error(errString)
	}

	AssetAsBytes, _ := APIstub.GetState(args[0])
	fmt.Println("returned asset:", AssetAsBytes)
	if AssetAsBytes == nil {
		return shim.Error(fmt.Sprintf("Customer record for PAN `%s` doesnt exists .... Register here !", args[0]))
	}
	customerData := customerAsset{}
	json.Unmarshal(AssetAsBytes, &customerData)
	return shim.Success([]byte(fmt.Sprintf("Record:" + customerData.Name + "," + strconv.Itoa(customerData.ID) + "," + customerData.Age + "," + customerData.Pan + "," + customerData.Sex)))
}

// issuePolicy
// Here the policy will be issue to customer based on the policy they prefer, one customer can have multiple different policy
// Param: customer ID and PAN
func (t *insuranceClaim) issuePolicy(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************************** issuePolicy")
	fmt.Println("no.of arguments:", args)

	// Check length of the arguments
	validateArgs, errString := validateArgs(args, 2)
	if !validateArgs {
		return shim.Error(errString)
	}

	// Check whether customer exists
	customerData := customerAsset{}
	customerAssetBytes, _ := APIstub.GetState(args[1])
	fmt.Println("returned customerAssetBytes:", customerAssetBytes)
	if customerAssetBytes == nil {
		return shim.Error(fmt.Sprintf("Customer record for PAN `%s` doesn't exist .... Kindly register first !", args[1]))
	}

	json.Unmarshal(customerAssetBytes, &customerData)
	id, _ := strconv.Atoi(args[0])
	if customerData.ID != id {
		fmt.Println("Customer record doesn't match", args[0], args[1])
		return shim.Error(fmt.Sprintf("Customer ID '%s' doesnt match with the given PAN '%s'", args[0], args[1]))
	}

	policyCurrentCounter := getCounter(APIstub, policyCounterNo)
	pn := fmt.Sprintf("AZ" + strconv.Itoa(policyCurrentCounter))
	fmt.Println("Issuing policy", pn)
	policyAssetNonBytes := policyAsset{
		PolicyNumber:     pn,
		CustomerID:       id,
		ClaimLastUpdated: time.Now(),
		ClaimStatus:      "NotClaimed",
	}

	policyAssetBytes, _ := json.Marshal(policyAssetNonBytes)
	errPut := APIstub.PutState(policyAssetNonBytes.PolicyNumber, policyAssetBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Not able to create policy for customer ID '%s' / PAN '%s'", args[0], args[1]))
	}

	policyCurrentCounter++
	counterBytes, _ := json.Marshal(counterNo{Counter: policyCurrentCounter})
	APIstub.PutState(policyCounterNo, counterBytes)
	fmt.Println("Successfully created policy ..... ", pn, args[0], args[1])
	return shim.Success([]byte(fmt.Sprintf("Successfully issued the policy, Policy ID: %s", pn)))
}

// claim
// Param: customer ID, policy ID
func (t *insuranceClaim) claim(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************************** claim")
	fmt.Println("no.of arguments:", args)

	// Check length of the arguments
	validateArgs, errString := validateArgs(args, 2)
	if !validateArgs {
		return shim.Error(errString)
	}

	policyAssetBytes, _ := APIstub.GetState(args[1])
	fmt.Println("returned policyAssetBytes:", policyAssetBytes)
	if policyAssetBytes == nil {
		return shim.Error(fmt.Sprintf("Policy ID: `%s` doesn't exists", args[1]))
	}

	policyData := policyAsset{}
	json.Unmarshal(policyAssetBytes, &policyData)
	id, _ := strconv.Atoi(args[0])
	fmt.Println("returned policyAssetBytes:", policyData)

	if policyData.CustomerID != id {
		return shim.Error(fmt.Sprintf("Policy Id `%s` doesnt exist for this customer `%s`", args[1], args[0]))
	}

	if policyData.ClaimStatus == "Claimed" {
		return shim.Error(fmt.Sprintf("Policy Id `%s`already claimed", args[1]))
	}

	policyData.ClaimLastUpdated = time.Now()
	policyData.ClaimStatus = "Claimed"

	policyDataBytes, _ := json.Marshal(policyData)
	errPut := APIstub.PutState(policyData.PolicyNumber, policyDataBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Update policy claim status failed .... Try again ! Policy ID: %s", args[1]))
	}

	return shim.Success([]byte("Policy successfully claimed"))
}

// queryAllAsset : Query all records
// param: None
func (t *insuranceClaim) queryAllAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************************** queryAllAsset")

	startKey := ""
	endKey := ""
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		// respValue := string(queryResponse.Value)
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- queryAllAssets:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

//=================== Private function ====================================
func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	fmt.Println("********************************** getCounter")
	fmt.Println("For:", AssetType)

	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := counterNo{}
	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Printf("Counter Current Value %d of Asset Type %s\n", counterAsset.Counter, AssetType)

	return counterAsset.Counter
}

func validateArgs(args []string, length int) (bool, string) {
	fmt.Println("********************************** validateArgs")
	fmt.Println("For:", args)

	if len(args) != length {
		return false, fmt.Sprintf("Incorrect number of arguments, Required %d arguments", length)
	}
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return false, fmt.Sprintf(string(i+1) + "st argument must be a non-empty string")
		}
	}
	return true, "Success"
}
