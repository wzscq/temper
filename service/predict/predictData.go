package predict

import (
	"temper/crv"
	"temper/common"
	"log"
	"os"
	"encoding/csv"
)	

type TemperRecItem struct{
	DeviceTypeID string `json:"temper_device_type_id"` 
	DeviceID string `json:"temper_device_id"` 
	SensorID string `json:"temper_sensor_id"` 
	Date string `json:"date"` 
	Time string `json:"time"` 
	Type string `json:"temper_type"` 
	Actual string `json:"temper_actual"`
	Predicted string `json:"temper_predicted"`
	ID string `json:"id"`
	Version string `json:"version"`
}

func GetPredictResultFileName(lastHisRecID string)(string){
	return lastHisRecID+"_result.csv"
}

func ReadResultRecsFromCSV(inFileName string)(*[]*TemperRecItem,error){
	//打开文件
	file,err:=os.Open(inFileName)
	if err!=nil {
		log.Println("ReadResultRecsFromCSV open file error",err.Error())
		return nil,err
	}
	defer file.Close()

	//创建csv reader
	reader:=csv.NewReader(file)
	//读取所有记录
	records,err:=reader.ReadAll()
	if err!=nil {
		log.Println("ReadResultRecsFromCSV read all error",err.Error())
		return nil,err
	}
	//解析记录
	var resultRecItems []*TemperRecItem
	for _,record:=range records {
		tempRecItem:=&TemperRecItem{
			DeviceTypeID:record[0],
			DeviceID:record[1],
			SensorID:record[2],
			Date:record[3],
			Time:record[4],
			Predicted:record[5],
		}
		resultRecItems=append(resultRecItems,tempRecItem)
	}

	return &resultRecItems,nil
}

func SaveHisRecsToCSV(outFileName string,hisRecItems *[]*TemperRecItem)(error){
	file,err:=os.Create(outFileName)
	if err!=nil {
		log.Println("SaveHisRecsToCSV create file error")
		return err
	}
	defer file.Close()

	writer:=csv.NewWriter(file)
	defer writer.Flush()

	for i:=len(*hisRecItems)-1;i>=0;i-- {
		hisRecItem:=(*hisRecItems)[i]
		err:=writer.Write([]string{
			hisRecItem.DeviceTypeID,
			hisRecItem.DeviceID,
			hisRecItem.SensorID,
			hisRecItem.Date,
			hisRecItem.Time,
			hisRecItem.Actual,
		})
		if err!=nil {
			log.Println("SaveHisRecsToCSV write error")
			return err
		}
	}

	return nil
}

func GetHisRecs(lastHisRecItem *TemperRecItem,histCount int,crvClient *crv.CRVClient,token string)(*[]*TemperRecItem,int){
	commonRep:=crv.CommonReq{
		ModelID:"temper_rec",
		Fields:&[]map[string]interface{}{
			map[string]interface{}{
				"field":"date",
			},
			map[string]interface{}{
				"field":"temper_actual",
			},
		},
		Filter:&map[string]interface{}{
			"temper_sensor_id":map[string]interface{}{
				"Op.eq":lastHisRecItem.SensorID,
			},
			"time":map[string]interface{}{
				"Op.eq":lastHisRecItem.Time,
			},
			"date":map[string]interface{}{
				"Op.lte":lastHisRecItem.Date,
			},
		},
		Pagination:&crv.Pagination{
			Current:1,
			PageSize:histCount,
		},
		Sorter:&[]crv.Sorter{
			crv.Sorter{
				Field:"date",
				Order:"desc",
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
		log.Println("GetLastHisRec can not be converted to map")
		return nil,common.ResultNoHisData
	}

	list,ok:=resultMap["list"]
	if !ok {
		log.Println("GetLastHisRec queryResult no list")
		return nil,common.ResultNoHisData
	}

	hisList,ok:=list.([]interface{})
	if !ok || len(hisList)<=0 {
		log.Println("GetLastHisRec queryResult no list")
		return nil,common.ResultNoHisData
	}

	var histRecItems []*TemperRecItem
	for _,item := range hisList {
		hisRec,ok:=item.(map[string]interface{})
		if !ok {
			log.Println("GetLastHisRec queryResult row 0 can not convert to map")
			return nil,common.ResultNoHisData
		}

		tempRecItem:=&TemperRecItem{
			DeviceTypeID:lastHisRecItem.DeviceTypeID,
			DeviceID:lastHisRecItem.DeviceID,
			SensorID:lastHisRecItem.SensorID,
			Time:lastHisRecItem.Time,
			Date:hisRec["date"].(string),
			Actual:hisRec["temper_actual"].(string),
		}

		histRecItems=append(histRecItems,tempRecItem)
	}

	return &histRecItems,common.ResultSuccess
}

func GetLastHisRec(lastHisRecID string,crvClient *crv.CRVClient,token string)(*TemperRecItem,int){
	commonRep:=crv.CommonReq{
		ModelID:"temper_rec",
		Fields:&[]map[string]interface{}{
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
		},
		Filter:&map[string]interface{}{
			"id":map[string]interface{}{
				"Op.eq":lastHisRecID,
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
		log.Println("GetLastHisRec can not be converted to map")
		return nil,common.ResultNoHisData
	}

	list,ok:=resultMap["list"]
	if !ok {
		log.Println("GetLastHisRec queryResult no list")
		return nil,common.ResultNoHisData
	}

	hisList,ok:=list.([]interface{})
	if !ok || len(hisList)<=0 {
		log.Println("GetLastHisRec queryResult no list")
		return nil,common.ResultNoHisData
	}

	//获取第一条记录
	hisRec,ok:=hisList[0].(map[string]interface{})
	if !ok {
		log.Println("GetLastHisRec queryResult row 0 can not convert to map")
		return nil,common.ResultNoHisData
	}

	tempRecItem:=&TemperRecItem{
		DeviceTypeID:hisRec["temper_device_type_id"].(string),
		DeviceID:hisRec["temper_device_id"].(string),
		SensorID:hisRec["temper_sensor_id"].(string),
		Date:hisRec["date"].(string),
		Time:hisRec["time"].(string),
	}

	return tempRecItem,common.ResultSuccess
}

func GetExistRecByDate(startRecItem,endRecItem *TemperRecItem,count int,crvClient *crv.CRVClient,token string)(*[]*TemperRecItem,int){
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
		},
		Filter:&map[string]interface{}{
			"temper_sensor_id":map[string]interface{}{
				"Op.eq":startRecItem.SensorID,
			},
			"time":map[string]interface{}{
				"Op.eq":startRecItem.Time,
			},
			"Op.and":[]map[string]interface{}{
				map[string]interface{}{
					"date":map[string]interface{}{
						"Op.lte":endRecItem.Date,
					},
				},
				map[string]interface{}{
					"date":map[string]interface{}{
						"Op.gte":startRecItem.Date,
					},
				},
			},
		},
		Pagination:&crv.Pagination{
			Current:1,
			PageSize:count,
		},
		Sorter:&[]crv.Sorter{
			crv.Sorter{
				Field:"date",
				Order:"asc",
			},
		},
	}

	rsp,errorCode:=crvClient.Query(&commonRep,token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	if rsp.Result==nil{
		return nil,common.ResultSuccess
	}

	//获取result中的list
	resultMap,ok:=rsp.Result.(map[string]interface{})
	if !ok {
		log.Println("GetExistRecByDate can not be converted to map")
		return nil,common.ResultSuccess
	}

	list,ok:=resultMap["list"]
	if !ok {
		log.Println("GetExistRecByDate queryResult no list")
		return nil,common.ResultSuccess
	}

	hisList,ok:=list.([]interface{})
	if !ok || len(hisList)<=0 {
		log.Println("GetExistRecByDate queryResult no list")
		return nil,common.ResultSuccess
	}

	var recItems []*TemperRecItem
	for _,item := range hisList {
		hisRec,ok:=item.(map[string]interface{})
		if !ok {
			log.Println("GetExistRecByDate queryResult row 0 can not convert to map")
			return nil,common.ResultSuccess
		}

		tempRecItem:=&TemperRecItem{
			ID:hisRec["id"].(string),
			Version:hisRec["version"].(string),
			Date:hisRec["date"].(string),
		}

		recItems=append(recItems,tempRecItem)
	}

	return &recItems,common.ResultSuccess
}

func SavePredictToDB(hisRecItems *[]*TemperRecItem,crvClient *crv.CRVClient,token string)(int){
	existRec,errorCode:=GetExistRecByDate((*hisRecItems)[0],(*hisRecItems)[len(*hisRecItems)-1],len(*hisRecItems),crvClient,token)
	if errorCode!=common.ResultSuccess {
		return errorCode
	}

	//创建数据保存列表
	var saveList []map[string]interface{}
	for _,hisRecItem:=range *hisRecItems {
		var extItem *TemperRecItem
		if existRec!=nil {
			for _,existItem:=range *existRec {
				if hisRecItem.Date==existItem.Date {
					extItem=existItem
					break
				}
			}
		}

		if extItem==nil {
			saveList=append(saveList,map[string]interface{}{
				"temper_device_type_id":hisRecItem.DeviceTypeID,
				"temper_device_id":hisRecItem.DeviceID,
				"temper_sensor_id":hisRecItem.SensorID,
				"date":hisRecItem.Date,
				"time":hisRecItem.Time,
				"temper_predicted":hisRecItem.Predicted,
				"_save_type":"create",
			})
		} else {
			saveList=append(saveList,map[string]interface{}{
				"id":extItem.ID,
				"version":extItem.Version,
				"temper_predicted":hisRecItem.Predicted,
				"_save_type":"update",
			})
		}
	}

	commonRep:=crv.CommonReq{
		ModelID:"temper_rec",
		List:&saveList,
	}

	_,errorCode=crvClient.Save(&commonRep,token)
	return errorCode
}