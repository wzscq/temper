package predict

import (
	"fmt"
	"testing"
	"temper/crv"
	"temper/common"
)

func _TestGetLastHisRec(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"predict_service",
	}

	tempRecItem,errorCode:=GetLastHisRec("1",crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetLastHisRec error")
	}

	fmt.Println(tempRecItem)
}

func TestGetHisRecs(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"predict_service",
	}

	lastHisRecItem:=&TemperRecItem{
		Date:"2018-06-05",
		Time:"00",
		SensorID:"1",
	}

	hisRecItems,errorCode:=GetHisRecs(lastHisRecItem,5,crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetLastHisRec error")
	}

	for _,tempRecItem:=range *hisRecItems {
		fmt.Println(*tempRecItem)
	}
}