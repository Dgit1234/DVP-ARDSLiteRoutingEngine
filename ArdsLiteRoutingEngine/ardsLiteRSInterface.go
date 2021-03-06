package main

import (
	"encoding/json"
	"fmt"
	"github.com/DuoSoftware/gorest"
)

type ArdsLiteRS struct {
	gorest.RestService `root:"/resourceselection/" consumes:"application/json" produces:"application/json"`
	getResource        gorest.EndPoint `method:"GET" path:"/getresource/{Company:int}/{Tenant:int}/{ResourceCount:int}/{SessionId:string}/{ServerType:string}/{RequestType:string}/{SelectionAlgo:string}/{HandlingAlgo:string}/{OtherInfo:string}" output:"string"`
}

func (ardsLiteRs ArdsLiteRS) GetResource(Company, Tenant, ResourceCount int, SessionId, ServerType, RequestType, SelectionAlgo, HandlingAlgo, OtherInfo string) string {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	fmt.Println("Company:", Company)
	fmt.Println("Tenant:", Tenant)
	fmt.Println("SessionId:", SessionId)
	fmt.Println("OtherInfo:", OtherInfo)

	byt := []byte(OtherInfo)
	var otherInfo string
	json.Unmarshal(byt, &otherInfo)

	requestKey := fmt.Sprintf("Request:%d:%d:%s", Company, Tenant, SessionId)

	if RedisCheckKeyExist(requestKey) {
		strReqObj := RedisGet(requestKey)
		fmt.Println(strReqObj)

		var reqObj Request
		json.Unmarshal([]byte(strReqObj), &reqObj)

		var result = SelectResources(Company, Tenant, ResourceCount, reqObj.LbIp, reqObj.LbPort, SessionId, ServerType, RequestType, SelectionAlgo, HandlingAlgo, otherInfo)
		return result
	}
	return "Session Invalied"

}

/*func GetRequestedResCount(ReqOtherInfo string) int {
	fmt.Println("ReqOtherInfo:", ReqOtherInfo)
	var requestedResCount = 1

	var resCount MultiResCount
	err := json.Unmarshal([]byte(ReqOtherInfo), &resCount)

	if err != nil {
		println(err)
		requestedResCount = 1
	} else {
		requestedResCount = resCount.ResourceCount
	}

	return requestedResCount
}*/

func SelectResources(Company, Tenant, ResourceCount int, ArdsLbIp, ArdsLbPort, SessionId, ServerType, RequestType, SelectionAlgo, HandlingAlgo, OtherInfo string) string {
	var selectionResult SelectionResult
	var handlingResult = ""
	switch SelectionAlgo {
	case "BASIC":
		selectionResult = BasicSelection(Company, Tenant, SessionId)
		break
	case "BASICTHRESHOLD":
		selectionResult = BasicThresholdSelection(Company, Tenant, SessionId)
		break
	case "WEIGHTBASE":
		selectionResult = WeightBaseSelection(Company, Tenant, SessionId)
		break
	case "LOCATIONBASE":
		selectionResult = LocationBaseSelection(Company, Tenant, SessionId)
		break
	default:
		break
	}

	switch HandlingAlgo {
	case "SINGLE":
		handlingResult = SingleResourceAlgo(ArdsLbIp, ArdsLbPort, ServerType, RequestType, SessionId, selectionResult, Company, Tenant)
	case "MULTIPLE":
		fmt.Println("ReqOtherInfo:", OtherInfo)
		resCount := ResourceCount
		fmt.Println("GetRequestedResCount:", resCount)
		handlingResult = MultipleHandling(ArdsLbIp, ArdsLbPort, ServerType, RequestType, SessionId, selectionResult, resCount, Company, Tenant)
	default:
		handlingResult = ""
	}

	return handlingResult
}
