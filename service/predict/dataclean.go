package predict

import (
	"temper/crv"
	"temper/common"
	"log"
	"os/exec"
)

func GetCleanResultFileName(recID string)(string){
	return recID+"_clean_result.csv"
}

func GetCleanFileName(recID string)(string){
	return recID+"_clean.csv"
}

func GetCleanItems(selectedRowKeys *[]string,crvClient *crv.CRVClient,token string)(*[]*TemperRecItem,int){
	commonRep:=crv.CommonReq{
		ModelID:"temper_rec",
		Fields:&[]map[string]interface{}{
			map[string]interface{}{
				"field":"id",
			},
			map[string]interface{}{
				"field":"date",
			},
			map[string]interface{}{
				"field":"time",
			},
			map[string]interface{}{
				"field":"temper_sensor_id",
			},
			map[string]interface{}{
				"field":"temper_device_id",
			},
			map[string]interface{}{
				"field":"temper_device_type_id",
			},
			map[string]interface{}{
				"field":"temper_original",
			},
		},
		Filter:&map[string]interface{}{
			"id":map[string]interface{}{
				"Op.in":selectedRowKeys,
			},
		},
		Pagination:&crv.Pagination{
			Current:1,
			PageSize:len(*selectedRowKeys),
		},
		Sorter:&[]crv.Sorter{
			crv.Sorter{
				Field:"id",
				Order:"asc",
			},
		},
	}

	rsp,errorCode:=crvClient.Query(&commonRep,token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	if rsp.Result==nil{
		return nil,common.ResultNoHisData
	}

	//获取result中的list
	resultMap,ok:=rsp.Result.(map[string]interface{})
	if !ok {
		return nil,common.ResultQueryRequestError
	}

	list,ok:=resultMap["list"].([]interface{})
	if !ok {
		return nil,common.ResultQueryRequestError
	}

	cleanItems:=make([]*TemperRecItem,0)
	for _,item:=range list {
		itemMap,ok:=item.(map[string]interface{})
		if !ok {
			return nil,common.ResultQueryRequestError
		}

		var orgValStr string
		orgVal,ok:=itemMap["temper_original"]
		if ok && orgVal!=nil {
			orgValStr=orgVal.(string)
		}

		cleanItem:=TemperRecItem{
			ID:itemMap["id"].(string),
			Time:itemMap["time"].(string),
			Date:itemMap["date"].(string),
			DeviceTypeID:itemMap["temper_device_type_id"].(string),
			DeviceID:itemMap["temper_device_id"].(string),
			SensorID:itemMap["temper_sensor_id"].(string),
			Actual:orgValStr,
		}

		cleanItems=append(cleanItems,&cleanItem)
	}

	return &cleanItems,common.ResultSuccess
}

func DataClean(cleanFileName,resultFileName string)(string){
	cmd := exec.Command("python3", "dataclean.py",cleanFileName,resultFileName)
  	// 设置执行命令时的工作目录
  	_, err := cmd.Output()
	if err != nil {
		log.Println("Predict exec dataclean.py error:",err)
		return ""
	}
	//log.Println(string(out))
  	return "0"
}

func IsSameDeviceSensor(recs *[]*TemperRecItem)(bool){
	senorid:=(*recs)[0].SensorID
	for _,item:= range *recs {
		if item.SensorID!=senorid {
			return false
		}
	}

	return true
}
