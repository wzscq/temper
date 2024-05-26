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

	selectedRowKeys:=[]string{"1","2","3","4","5"}

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
	err:=SaveRecsToCSV(cleanFileName,cleanItems)
	if err!=nil {
		t.Error("SaveRecsToCSV error")
	}

	
}