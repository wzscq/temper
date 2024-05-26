package predict

import (
	"temper/crv"
	"temper/common"
	"log"
	"os/exec"
	"os"
	"encoding/csv"
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
				"field":"version",
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
				Field:"date",
				Order:"asc",
			},
			crv.Sorter{
				Field:"time",
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
			Version:itemMap["version"].(string),
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

func SaveCleanRecsToDB(hisRecItems *[]*TemperRecItem,crvClient *crv.CRVClient,token string)(int){
	var saveList []map[string]interface{}
	for _,hisRecItem:=range *hisRecItems {

		saveList=append(saveList,map[string]interface{}{
			"id":hisRecItem.ID,
			"version":hisRecItem.Version,
			"temper_actual":hisRecItem.Predicted,
			"_save_type":"update",
		})
	}

	commonRep:=crv.CommonReq{
		ModelID:"temper_rec",
		List:&saveList,
	}

	_,errorCode:=crvClient.Save(&commonRep,token)
	return errorCode
}

func SaveCleanRecsToCSV(outFileName string,hisRecItems *[]*TemperRecItem)(error){
	file,err:=os.Create(outFileName)
	if err!=nil {
		log.Println("SaveCleanRecsToCSV create file error")
		return err
	}
	defer file.Close()

	writer:=csv.NewWriter(file)
	defer writer.Flush()

	for i:=len(*hisRecItems)-1;i>=0;i-- {
		hisRecItem:=(*hisRecItems)[i]
		err:=writer.Write([]string{
			hisRecItem.ID,
			hisRecItem.Version,
			hisRecItem.DeviceTypeID,
			hisRecItem.DeviceID,
			hisRecItem.SensorID,
			hisRecItem.Date,
			hisRecItem.Time,
			hisRecItem.Actual,
		})
		if err!=nil {
			log.Println("SaveCleanRecsToCSV write error")
			return err
		}
	}

	return nil
}

func ReadCleanResultRecsFromCSV(inFileName string)(*[]*TemperRecItem,error){
	//打开文件
	file,err:=os.Open(inFileName)
	if err!=nil {
		log.Println("ReadCleanResultRecsFromCSV open file error",err.Error())
		return nil,err
	}
	defer file.Close()

	//创建csv reader
	reader:=csv.NewReader(file)
	//读取所有记录
	records,err:=reader.ReadAll()
	if err!=nil {
		log.Println("ReadCleanResultRecsFromCSV read all error",err.Error())
		return nil,err
	}
	//解析记录
	var resultRecItems []*TemperRecItem
	for _,record:=range records {
		tempRecItem:=&TemperRecItem{
			ID:record[0],
			Version:record[1],
			DeviceTypeID:record[2],
			DeviceID:record[3],
			SensorID:record[4],
			Date:record[5],
			Time:record[6],
			Predicted:record[7],
		}
		resultRecItems=append(resultRecItems,tempRecItem)
	}

	return &resultRecItems,nil
}
