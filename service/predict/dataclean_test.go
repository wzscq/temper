package predict

import (
	"fmt"
	"testing"
	"temper/crv"
	"temper/common"
)

func TestGetCleanItems(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"predict_service",
	}

	selectedRowKeys:=[]string{"1","177","353","529","705","881","1057","1233","1409","1585","1761","1937","2113","2289","2465","2641","2817","2993","3169","3345","3521","3697","3873","4049"}

	cleanItems,errorCode:=GetCleanItems(&selectedRowKeys,crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("GetCleanItems error")
	}
	fmt.Println(cleanItems)

	if len(*cleanItems)<=0 {
		t.Error("no recs")
	}

	if IsSameDeviceSensor(cleanItems) == false {
		t.Error("not the same sensor")
	}

	lastRecID:=(*cleanItems)[0].ID
	//获取训练数据输出文件名
	cleanFileName:=GetCleanFileName(lastRecID)

	//保持历史记录到CSV文件
	err:=SaveCleanRecsToCSV(cleanFileName,cleanItems)
	if err!=nil {
		t.Error("SaveCleanRecsToCSV error")
	}

	//加载数据
	resultRecs,err:=ReadCleanResultRecsFromCSV(cleanFileName)
	if err!=nil {
		t.Error("ReadCleanResultRecsFromCSV error")
	}

	//保存数据到数据库
	errorCode=SaveCleanRecsToDB(resultRecs,crvClient,"predict_service")
	if errorCode!=common.ResultSuccess {
		t.Error("SaveCleanRecsToDB error")
	}
}