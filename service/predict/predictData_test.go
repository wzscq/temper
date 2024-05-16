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

func _TestGetHisRecs(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"predict_service",
	}

	lastHisRecItem,errorCode:=GetLastHisRec("342",crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetLastHisRec error")
	}

	hisRecItems,errorCode:=GetHisRecs(lastHisRecItem,100,crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetLastHisRec error")
	}

	for _,tempRecItem:=range *hisRecItems {
		fmt.Println(*tempRecItem)
	}

	fileName:=GetHisFileName("1")
	err:=SaveHisRecsToCSV(fileName,hisRecItems)
	if err!=nil {
		t.Error("SaveHisRecsToCSV error")
	}
}

func TestReadResultRecsFromCSV(t *testing.T) {
	fileName:=GetResultFileName("342")
	resultRecs,err:=ReadResultRecsFromCSV(fileName)
	if err!=nil {
		t.Error("ReadResultRecsFromCSV error")
	}

	fmt.Println("ReadResultRecsFromCSV:")
	for _,resultRec:=range *resultRecs {
		fmt.Println(*resultRec)
	}

	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"predict_service",
	}

	/*existRec,errorCode:=GetExistRecByDate((*resultRecs)[0],(*resultRecs)[len(*resultRecs)-1],len(*resultRecs),crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetExistRecByDate error")
	}
	for _,tempRecItem:=range *existRec {
		fmt.Println(*tempRecItem)
	}*/

	if len(*resultRecs)>0 {
		errCode:=SavePredictToDB(resultRecs,crvClient,"predict_service")
		if errCode!=common.ResultSuccess {
			t.Error("GetExistRecByDate error")
		}
		fmt.Println("GetExistRecByDate:")
	}
}

func _TestGetExistRecByDate(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"predict_service",
	}

	startItem:=&TemperRecItem{
		Time:"01",
		SensorID:"1",
		Date:"2018-04-11 00:00:00",
	}

	endItem:=&TemperRecItem{
		Time:"01",
		SensorID:"1",
		Date:"2018-04-17 00:00:00",
	}

	hisRecItems,errorCode:=GetExistRecByDate(startItem,endItem,7,crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetExistRecByDate error")
	}

	fmt.Println("GetExistRecByDate:")
	for _,tempRecItem:=range *hisRecItems {
		fmt.Println(*tempRecItem)
	}

}