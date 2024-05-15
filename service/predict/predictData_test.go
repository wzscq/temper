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

	lastHisRecItem,errorCode:=GetLastHisRec("4223",crvClient,"predict_service")
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
	fileName:=GetResultFileName("1")
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

	if len(*resultRecs)>0 {
		errCode:=SavePredictToDB(resultRecs,crvClient,"predict_service")
		if errCode!=common.ResultSuccess {
			t.Error("GetExistRecByDate error")
		}
		fmt.Println("GetExistRecByDate:")
	}
}