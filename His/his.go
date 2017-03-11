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
	"strings"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
// myChaincode example simple Chaincode implementation
type myChaincode struct {
}
var itemSp = "\n"
var listSp = "!@#$"
var minSp = "^&*"
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(myChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
// Init resets all the things
func (t *myChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}
// Invoke is our entry point to invoke a chaincode function
func (t *myChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "register":
		return t.register(stub, args)
	case "add":
		return t.add(stub, args)
	case "updateAccessInfo":
		return t.updateAccessInfo(stub, args)
	case "updatePermission":
		return t.updatePermission(stub, args)
	case "confirmSummary":
		return t.confirmSummary(stub, args)
	case "setCondition":
		return t.setCondition(stub, args)
	case "payForRecord":
		return t.payForRecord(stub, args)
	default:
		return nil, errors.New("Unsupported operation")
	}
}
// Query is our entry point for queries
func (t *myChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "havePermission":
		return t.havePermission(stub, args)
	case "getRecord":
		return t.getRecord(stub, args)
	case "getSummary":
		return t.getSummary(stub, args)
	default:
		return nil, errors.New("Unsupported operation")
	}
}
func (t *myChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 2{
		return nil, errors.New("register operation must include at last 1 argument : userid, balance")
	}
	userId := args[0]
	balance := args[1]
	//check the existence of the user
	value, err := stub.GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("register operation failed. Error accessing state (check the existence of the user)")
	}
	if value != nil {
		return nil, fmt.Errorf("exited user")
	}
	defaultValue := "0" + itemSp + "0" + itemSp + "0" + minSp + "0" + itemSp + balance
	err = stub.PutState(userId, []byte(defaultValue))
	if err != nil {
		return nil, fmt.Errorf("register operation failed. Error putting starte: " + err.Error())
	}
	return nil, nil
}
func (t *myChaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	//(recordId, ownerId, providerId, accessInfo, permission)
	if len(args) < 8 {
		return nil, errors.New("add operation must include at last 5 arguments: （recordId, ownerId, providerId, accessInfo, permission[dataItem,  query, hash, deadline]) ")
	}
	//get the args
	recordId := args[0]
	ownerId := args[1]
	providerId := args[2]
	accessInfo := args[3]
	dataItem := args[4]
	query := args[5]
	hash := args[6]
	deadline := args[7]
	permission := ownerId + minSp + dataItem + minSp + query + minSp + hash + minSp + deadline
	//check the existence of the recordId
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("add operation failed. Error accessing state (check the existence of the recordId")
	}
	if value != nil {
		return nil, fmt.Errorf("existed recordId")
	}
	//write the record
	newValue := ownerId + itemSp + providerId + itemSp + accessInfo + itemSp + permission + itemSp + "0" + listSp + "0"
	err = stub.PutState(recordId, []byte(newValue))
	if err!= nil {
		return nil, fmt.Errorf("add operation failed.  Error while writing record: " +  err.Error())
	}
	//update the summary of the owner
	value, err = stub.GetState(ownerId)
	if err != nil {
		return nil, fmt.Errorf("add operation failed. Error accessing state (check the existence of the ownerId")
	}
	if value == nil {
		return nil, fmt.Errorf("not exist woner")
	}
	listValue := strings.Split(string(value), itemSp)
	count := listValue[1]
	recordList := listValue[2]
	newRecordR := recordId + minSp + "1"
	balance := listValue[3]
	intCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("add operation failed. strconv err")
	}
	count = strconv.Itoa(intCount + 1)
	newValue =  "1" + itemSp + count + itemSp + recordList + listSp + newRecordR + itemSp + balance
	err = stub.PutState(ownerId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("add operation failed. Error while update the Onwer's suammary: " + err.Error())
	}
	//update the summary of the provider
	value, err = stub.GetState(providerId)
	if err != nil {
		return nil, fmt.Errorf("add operation failed. Error accessing state (check the existence of the ownerId")
	}
	if value == nil {
		return nil, fmt.Errorf("not exist provider")
	}
	listValue = strings.Split(string(value), itemSp)
	flag := listValue[0]
	count = listValue[1]
	recordList = listValue[2]
	balance = listValue[3]
	newRecordR = recordId + minSp + "1"
	intCount, err = strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("add operation failed. strconv err")
	}
	count= strconv.Itoa(intCount + 1)
	newValue = flag + itemSp + count + itemSp + recordList + listSp + newRecordR + itemSp + balance
	err = stub.PutState(providerId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("add operation failed. Error while update the provider's suammary: " + err.Error())
	}
	return nil,nil
}
func (t *myChaincode) updateAccessInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 3{
		return nil, fmt.Errorf("create operation must include at last 3 arguments: （recordId, providerId, accessInfo) ")
	}
	//get the args
	recordId := args[0]
	_providerId := args[1]
	accessInfo := args[2]
	//get the record info
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("updateAccessInfo operation failed. Error while accessing state :" + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("the record not exists")
	}
	listValue := strings.Split(string(value), itemSp)
	ownerId := listValue[0]
	providerId := listValue[1]
	permission := listValue[3]
	condition := listValue[4]
	//check the provider
	if ( _providerId != providerId ){
		return nil, fmt.Errorf("don't have the right to update the accessinfo")
	}
	newValue := ownerId + itemSp + providerId + itemSp + accessInfo + itemSp + permission + itemSp + condition
	err = stub.PutState(recordId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("updateAccessInfo operation failed. Error while updating state :" + err.Error())
	}
	return nil, nil
}
func (t *myChaincode) updatePermission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 7{
		return nil, fmt.Errorf("updatePermission operation must include at last 7 arguments: (recordId, ownerId, [consumer, dataItem,  query, hash, deadline]) ")
	}
	//get the args
	recordId := args[0]
	_ownerId := args[1]
	consumer := args[2]
	dataItem := args[3]
	query := args[4]
	hash := args[5]
	deadline := args[6]
	//read the record
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("updatePermission operation failed. Error while getting the record: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this record")
	}
	listValue := strings.Split(string(value), itemSp)
	ownerId := listValue[0]
	providerId := listValue[1]
	accessInfo := listValue[2]
	condition := listValue[4]
	//check the right
	if _ownerId != ownerId {
		return nil, fmt.Errorf("don't have the right to do this")
	}
	permission := listValue[3]
	//maybe lookup the list and find out the existent permission, then delete the old one
	permission = permission + listSp + consumer + minSp + dataItem + minSp + query + minSp + hash + minSp + deadline
	newValue := ownerId + itemSp + providerId + itemSp + accessInfo + itemSp + permission + itemSp + condition

	//update the record' permission
	err = stub.PutState(recordId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("updatePermission operation failed. Error while updating record: " + err.Error())
	}
	//update the consumer's summary
	value, err = stub.GetState(consumer)
	if err != nil {
		return nil, fmt.Errorf("updatePermission operation failed. Error while getting consumer's sumamary: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("the consumer not exists")
	}
	listValue = strings.Split(string(value), itemSp)
	count := listValue[1]
	recordList := listValue[2]
	balance := listValue[3]
	//maybe lookup the list and find out the existent recorid, then delete the old one
	newRecordR := recordId + minSp + "1"
	intCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("updatePermission operation failed. strconv err")
	}
	count = strconv.Itoa(intCount + 1)

	newvalue := "1" + itemSp + count + itemSp + recordList + listSp + newRecordR + itemSp + balance
	err = stub.PutState(consumer, []byte(newvalue))
	if err != nil {
		return nil, fmt.Errorf("updatePermission operation failed. Error while update the provider's suammary: " + err.Error())
	}
	return nil, nil
}
func (t *myChaincode) confirmSummary(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 3{
		return nil, fmt.Errorf("confirmSummary operation must include at last 3 arguments: （ownerId, recordId, op) ")
	}
	//get the args
	ownerId := args[0]
	recordId := args[1]
	op := args[2]
	if op != "ac" && op != "re"{
		return nil, fmt.Errorf("bad format of the op")
	}
	//get the owner's summary
	value, err := stub.GetState(ownerId)
	if err != nil {
		return nil, fmt.Errorf("confirmSummary operation failed. Error while getting owner's sumamary: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("the owner not exists")
	}
	listValue := strings.Split(string(value), itemSp)
	flag := listValue[0]
	count := listValue[1]
	recordList := listValue[2]
	balance := listValue[3]
	//check the flag of the summary
	if flag == "0" {
		return nil, fmt.Errorf("there is nothing changed:(")
	}
	//update the status of the recordId
	listRecord := strings.Split(string(recordList), listSp)
	newRecordList := ""
	existence := false
	noNew := true
	for _, record := range listRecord{
		if record == "" {
			continue
		}
		rcd := strings.Split(record, minSp)
		//find the one
		if rcd[0] == recordId{
			existence = true
			if op == "ac"{
				newRecord := recordId + minSp + "0"
				newRecordList = newRecordList + listSp + newRecord
				//update the provider's summary?????
			} else {
				//update the count
				intCount, err := strconv.Atoi(count)
				if err != nil {
					return nil, fmt.Errorf("confirmSummary operation failed. strconv err")
				}
				count = strconv.Itoa( intCount - 1)
				//update the provider's summary???

			}
		} else {
			newRecordList = newRecordList + listSp + record
			if rcd[1] == "1"{
				noNew = false
			}
		}
	}
	if !existence {
		return nil, fmt.Errorf("the recordId not exists")
	}
	if noNew {
		flag = "0"
	}
	//do the update
	newValue := flag + itemSp + count + itemSp + newRecordList + itemSp + balance
	err = stub.PutState(ownerId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("confirmSummary operation failed. Error while update the provider's suammary: " + err.Error())
	}
	return nil, nil
}
func (t *myChaincode) setCondition(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 7{
		return nil, fmt.Errorf("confirmSummary operation must include at last 7 arguments: （userId, recordId, price, condition[dataItem,  query, hash, deadline])")
	}
	//get the args
	userId := args[0]
	recordId := args[1]
	price := args[2]
	dataItem := args[3]
	query := args[4]
	hash := args[5]
	deadline := args[6]
	//check the format of the price
	_, err := strconv.Atoi(price)
	if err != nil {
		return nil, fmt.Errorf("bad format of the price")
	}
	//get the owner's summary
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("setCondition operation failed. Error while getting owner's sumamary: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("the recordId not exists")
	}
	listValue := strings.Split(string(value), itemSp)
	ownerId := listValue[0]
	providerId := listValue[1]
	accessInfo := listValue[2]
	permission := listValue[3]
	//condition := listValue[4]
	//check the right
	if userId != ownerId {
		return nil, fmt.Errorf("don't have the right")
	}
	permissionWithoutConsumer := dataItem + minSp + query + minSp + hash + minSp + deadline
	condition := price + listSp + permissionWithoutConsumer
	newValue := ownerId + itemSp + providerId + itemSp + accessInfo + itemSp + permission + itemSp + condition
	err = stub.PutState(recordId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("setCondition operation failed. Error while putting condition : " + err.Error())
	}
	return nil, nil
}
func (t *myChaincode) payForRecord(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 2{
		return nil, fmt.Errorf("confirmSummary operation must include at last 2 arguments: （userId, recordId)")
	}
	//get the args
	userId := args[0]
	recordId := args[1]
	//get the info of the record
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while getting user's record : " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this record")
	}
	listValue := strings.Split(string(value), itemSp)
	ownerId := listValue[0]
	providerId := listValue[1]
	accessInfo := listValue[2]
	permission := listValue[3]
	condition := listValue[4]
	listCondition := strings.Split(condition, listSp)
	price := listCondition[0]
	permissionWithoutConsumer := listCondition[1]
	if permissionWithoutConsumer == "0"{
		return nil, fmt.Errorf("the woner don't want to share this record")
	}
	//get the info of the user
	value, err = stub.GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while getting user's summary : " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this user")
	}
	listValue = strings.Split(string(value), itemSp)
	flag := listValue[0]
	count := listValue[1]
	recordList := listValue[2]
	balance := listValue[3]
	//check the balance
	intBalance, err := strconv.Atoi(balance)
	if err != nil{
		return nil, fmt.Errorf("payForRecord operation failed. Error while convert the string to int: " + err.Error())
	}
	intPrice, err := strconv.Atoi(price)
	if err != nil{
		return nil, fmt.Errorf("payForRecord operation failed. Error while convert the string to int: " + err.Error())
	}
	if intBalance < intPrice {
		return nil, fmt.Errorf("don't have enough balance")
	}
	intBalance = intBalance - intPrice
	//update the user's suammary
	intCount, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while convert the string to int: " + err.Error())
	}
	intCount = intCount +1
	count = strconv.Itoa(intCount)
	balance = strconv.Itoa(intBalance)
	recordList = recordList + listSp + recordId + minSp + "0"
	newValue := flag + itemSp + count + itemSp + recordList + itemSp + balance
	err = stub.PutState(userId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while putting user's summary : " + err.Error())
	}
	//get the owner's summary
	value, err = stub.GetState(ownerId)
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while getting onwer's summary : " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this owner")
	}
	listValue = strings.Split(string(value), itemSp)
	flag = listValue[0]
	count = listValue[1]
	recordList = listValue[2]
	balance = listValue[3]
	//check the balance
	intBalance, err = strconv.Atoi(balance)
	if err != nil{
		return nil, fmt.Errorf("payForRecord operation failed. Error while convert the string to int: " + err.Error())
	}
	intBalance = intBalance + intPrice
	//update the owner's suammary
	balance = strconv.Itoa(intBalance)
	newValue = flag + itemSp + count + itemSp + recordList + itemSp + balance
	err = stub.PutState(ownerId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while putting owner's summary : " + err.Error())
	}
	//update the record's permission list
	permission = permission + listSp + userId + minSp + permissionWithoutConsumer
	newValue = ownerId + itemSp + providerId + itemSp + accessInfo + itemSp + permission + itemSp + condition
	err = stub.PutState(recordId, []byte(newValue))
	if err != nil {
		return nil, fmt.Errorf("payForRecord operation failed. Error while putting the record : " + err.Error())
	}
	return nil, nil
}
func (t *myChaincode) havePermission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 3{
		return nil, errors.New("havePermission operation must include at last 2 arguments: (userId, recordId, query) ")
	}
	//get the args
	userId := args[0]
	recordId := args[1]
	query := args[2]
	ret := "False"
	//read the record
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("havePermission operation failed. Error while getting the record: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this record")
	}
	listValue := strings.Split(string(value), itemSp)
	ownerId := listValue[0]
	providerId := listValue[1]
	//accessInfo := listValue[2]
	permission := listValue[3]
	//condition := listValue[4]
	//if the user is provider or user
	if userId == ownerId || userId == providerId {
		return []byte("True"), nil
	}
	listPermission := strings.Split(permission, listSp)
	for _, p := range listPermission {
		listP := strings.Split(p, minSp)
		if listP[0] == userId && listP[2] == query{
			ret = "True"
			break
		}
	}
	return []byte(ret), nil
}
func (t *myChaincode) getRecord(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 2{
		return nil, errors.New("getRecord operation must include at last 2 arguments: (userId, recordId)")
	}
	//get the args
	userId := args[0]
	recordId := args[1]
	//query := args[2]
	//read the record
	value, err := stub.GetState(recordId)
	if err != nil {
		return nil, fmt.Errorf("getRecord operation failed. Error while getting the record: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this record")
	}
	listValue := strings.Split(string(value), itemSp)
	ownerId := listValue[0]
	providerId := listValue[1]
	accessInfo := listValue[2]
	permission := listValue[3]
	//if the user is provider or user
	if userId == ownerId || userId == providerId {
		return value, nil
	}
	ret := ownerId + itemSp + providerId + itemSp + accessInfo
	mypermission := ""
	listPermission := strings.Split(permission, listSp)
	for _, p := range listPermission {
		listP := strings.Split(p, minSp)
		if listP[0] == userId {
			mypermission = p
			break
		}
	}
	if mypermission == "" {
		return nil, fmt.Errorf("don't have the right to get this record")
	}
	ret = ret + itemSp + mypermission
	return []byte(ret), nil
}
func (t *myChaincode) getSummary(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) < 1{
		return nil, errors.New("getSummary operation must include at last 1 argument: (userId) ")
	}
	userId := args[0]
	//read the summary
	value, err := stub.GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("getSummary operation failed. Error while getting the record: " + err.Error())
	}
	if value == nil {
		return nil, fmt.Errorf("don't have this user")
	}
	return value, nil
}
