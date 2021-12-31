package main

import (
  "fmt"
  "encoding/json"
  "log"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
  contractapi.Contract
}

/**Credit describes the basic details of a response credit
 * the struct fields are inserted in alphabetic order to ensure that 
 * the order of fields in json format are the same when called and 
 * returned by different programming languages
 */
/**Hyperledger Fabric supports blockchain with world state. The world state identifies the current state of each object, while
 * the blockchain records the changes of the world state. The data structure of the object correspond to the world state, while
 * the transactions are recorded in the blockchain.
 */
/**The data structure of credit 
 * ID           "000001"        string
 * IssueDate    "2021-12-30"    string
 * Issuer       "VPP Operator"  string
 * Owner        "VPP Operator" "DER1" ... string
 * State        "Issued" "Trading" "Redeemed" string
 * Value        5  int (this represents the value of the response at current stage)
 */

/**Transaction type: issue, trade, redeem
 * Issue: generation of the credit, endorsed by the issuer (which is also the initial owner)
 *    the issuing transaction should include the following information:
 *        Type            = "Issue"
 *        Issuer          = "VPP Operator"
 *        ID              = "000001"
 *        Issue time      = "2021-12-30 12:00"
 *        Value           = 5
 * Trade: change of owner of the credit, endorsed by both the original and the new owners
 *    the trading transaction should include the following information:
 *        Type            = "Trade"
 *        Issuer          = "VPP Operator"
 *        ID              = "000001"
 *        Current owner   = "BESS1"
 *        New owner       = "BESS2"
 *        Trade time      = "2021-12-31 10:00"
 *        Price           = 4
 * Redeem: expiration of the credit, endorsed by the current owner and the issuer
 *    the redeem transaction should include the following information:
 *        Type            = "Redeem"
 *        Issuer          = "VPP Operator"
 *        ID              = "000001"
 *        Current owner   = "BESS1"
 *        Redeem time     = "2021-12-31 15:00"
 */

type Credit struct {
  ID             string `json:"ID"`
  IssueDate      string `json:"IssueDate"`
  Owner          string `json:"Owner"`
}

// InitLedger adds a base set of assets to the ledger
/*
here InitLedger is defined as a method of an interface
s is the name of the structure variable, * means we are invoking the pointer of the variable, SmartContract is the name of the
struct we are calling. ctx is the input and contractapi.TransactionContextInterface is the type of the input, and error is the
only output
*/
/* Note that when using the contract api, each chaincode function that is called is passed a transaction context “ctx”, 
from which you can get the chaincode stub (GetStub() ), which has functions to access the ledger (e.g. GetState() ) and make 
requests to update the ledger (e.g. PutState() ). 
*/
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
  // struct initiialization using key-value pair, [] means that assets is an array
  credits := []Credit{
    {ID: "000001", IssueDate: "2021-12-25", Owner: "BESS1"}, 
    {ID: "000002", IssueDate: "2021-12-25", Owner: "BESS1"}, 
    {ID: "000003", IssueDate: "2021-12-25", Owner: "BESS1"}, 
  }

  // range returns the pointer and the content of each element in the array
  for _, credit := range credits {
    creditJSON, err := json.Marshal(credit)
    // nil means that the error message is empty
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(credit.ID, creditJSON)
    if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
    }
  }

  return nil
}
// * keep the "nil" after return
// * noted that there is always an output for the error message, when there is no error, nil should be returned to that position
// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) IssueCredit(ctx contractapi.TransactionContextInterface, id string, issuedate string, owner string) error {
  // query the smart contract to see whether the asset exist
  exists, err := s.CreditExists(ctx, id)
  if err != nil {
    return err
  }
  if exists {
    return fmt.Errorf("the credit %s already exists", id)
  }

  credit := Credit{
    ID:             id,
    IssueDate:      issuedate,
    Owner:          owner,
  }
  creditJSON, err := json.Marshal(credit)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, creditJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadCredit(ctx contractapi.TransactionContextInterface, id string) (*Credit, error) {
  creditJSON, err := ctx.GetStub().GetState(id)
  if err != nil {
    return nil, fmt.Errorf("failed to read from world state: %v", err)
  }
  if creditJSON == nil {
    return nil, fmt.Errorf("the credit %s does not exist", id)
  }

  var credit Credit
  err = json.Unmarshal(creditJSON, &credit)
  if err != nil {
    return nil, err
  }

  return &credit, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateCredit(ctx contractapi.TransactionContextInterface, id string, issuedate string, owner string) error {
  exists, err := s.CreditExists(ctx, id)
  if err != nil {
    return err
  }
  if !exists {
    return fmt.Errorf("the credit %s does not exist", id)
  }

  // overwriting original asset with new asset
  credit := Credit{
    ID:             id,
    IssueDate:      issuedate,
    Owner:          owner,
  }
  creditJSON, err := json.Marshal(credit)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, creditJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteCredit(ctx contractapi.TransactionContextInterface, id string) error {
  exists, err := s.CreditExists(ctx, id)
  if err != nil {
    return err
  }
  if !exists {
    return fmt.Errorf("the credit %s does not exist", id)
  }

  return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) CreditExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
  creditJSON, err := ctx.GetStub().GetState(id)
  if err != nil {
    return false, fmt.Errorf("failed to read from world state: %v", err)
  }

  return creditJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferCredit(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
  credit, err := s.ReadCredit(ctx, id)
  if err != nil {
    return err
  }

  credit.Owner = newOwner
  creditJSON, err := json.Marshal(credit)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, creditJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllCredits(ctx contractapi.TransactionContextInterface) ([]*Credit, error) {
  // range query with empty string for startKey and endKey does an
  // open-ended query of all assets in the chaincode namespace.
  resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
  if err != nil {
    return nil, err
  }
  defer resultsIterator.Close()

  var credits []*Credit
  for resultsIterator.HasNext() {
    queryResponse, err := resultsIterator.Next()
    if err != nil {
      return nil, err
    }

    var credit Credit
    err = json.Unmarshal(queryResponse.Value, &credit)
    if err != nil {
      return nil, err
    }
    credits = append(credits, &credit)
  }

  return credits, nil
}

func main() {
  creditChaincode, err := contractapi.NewChaincode(&SmartContract{})
  if err != nil {
    log.Panicf("Error creating response-credit chaincode: %v", err)
  }

  if err := creditChaincode.Start(); err != nil {
    log.Panicf("Error starting response-credit chaincode: %v", err)
  }
}