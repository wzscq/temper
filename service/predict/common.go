package predict

import (
	"os"
	"log"
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

func SaveRecsToCSV(outFileName string,hisRecItems *[]*TemperRecItem)(error){
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
