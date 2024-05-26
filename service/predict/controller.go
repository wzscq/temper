package predict

import (
	"temper/crv"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"temper/common"
)

type PredictController struct {
	CRVClient *crv.CRVClient
	PredictHisCount int
	DataCleanHisCount int
	DataCleanFollowingCount int
}

func (pc *PredictController) Bind(router *gin.Engine) {
	router.POST("/predict", pc.Predict)
	router.POST("/dataclean", pc.DataClean)
}

func (pc *PredictController) DataClean(c *gin.Context) {
	log.Println("start PredictController DataClean")

	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end PredictController DataClean with error")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end PredictController DataClean with error")
		return
  	}	

	if rep.SelectedRowKeys==nil || len(*rep.SelectedRowKeys)==0 {
		log.Println("end PredictController DataClean with error:SelectedRowKeys is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	cleanItems,errorCode:=GetCleanItems(rep.SelectedRowKeys,pc.CRVClient,header.Token)
	if errorCode!=common.ResultSuccess {
		log.Println("end PredictController DataClean with error:errorCode=",errorCode)
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	if len(*cleanItems)<pc.DataCleanHisCount {
		rsp:=common.CreateResponse(common.CreateError(common.ResultFollowingRecNotEnough,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//判断数据是否来自同一个测温点
	if IsSameDeviceSensor(cleanItems) == false {
		rsp:=common.CreateResponse(common.CreateError(common.ResultRecNotTheSameSensor,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	lastRecID:=(*cleanItems)[0].ID
	//获取训练数据输出文件名
	cleanFileName:=GetCleanFileName(lastRecID)

	//保存数据到CSV文件
	err:=SaveRecsToCSV(cleanFileName,cleanItems)
	if err!=nil {
		log.Println("end PredictController DataClean with error:SaveHisRecsToCSV error")
		rsp:=common.CreateResponse(common.CreateError(common.ResultSaveHisrecToCSVError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//调用预测服务
	resultFileName:=GetCleanResultFileName(lastRecID)
	result:=DataClean(cleanFileName,resultFileName)
	if result!="0" {
		log.Println("end PredictController DataClean with error:Predict error")
		params:=map[string]interface{}{
			"message":result,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultPredictError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//读取结果数据
	resultRecs,err:=ReadResultRecsFromCSV(resultFileName)
	if err!=nil {
		log.Println("end PredictController DataClean with error:ReadResultRecsFromCSV error")
		rsp:=common.CreateResponse(common.CreateError(common.ResultReadResultRecsFromCSVError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//保存结果数据到数据库
	errorCode=SaveRecsToDB(resultRecs,pc.CRVClient,header.Token)
	if errorCode!=common.ResultSuccess {
		log.Println("end PredictController DataClean with error:SavePredictToDB error")
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (pc *PredictController) Predict(c *gin.Context) {
	log.Println("start PredictController Predict")

	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end PredictController Predict with error")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("end PredictController Predict with error")
		return
  	}	

	if rep.SelectedRowKeys==nil || len(*rep.SelectedRowKeys)==0 {
		log.Println("end PredictController Predict with error:SelectedRowKeys is empty")
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//lastHisRecID:=(*rep.SelectedRowKeys)[0]
	for _,lastHisRecID := range *rep.SelectedRowKeys {
		log.Println("lastHisRecID:",lastHisRecID)
		//获取最后一条历史记录
		lastHisRecItem,errorCode:=GetLastHisRec(lastHisRecID,pc.CRVClient,header.Token)
		if errorCode!=common.ResultSuccess {
			log.Println("end PredictController Predict with error:errorCode=",errorCode)
			rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		//获取用于预测的历史记录
		hisRecItems,errorCode:=GetHisRecs(lastHisRecItem,pc.PredictHisCount,pc.CRVClient,header.Token)
		if errorCode!=common.ResultSuccess {
			log.Println("end PredictController Predict with error:errorCode=",errorCode)
			rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		//如果历史记录不足，返回错误
		if len(*hisRecItems)<pc.PredictHisCount {
			log.Println("end PredictController Predict with error:ResultNoHisData")
			rsp:=common.CreateResponse(common.CreateError(common.ResultHisRecNotEnough,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		//获取训练数据输出文件名
		hisFileName:=GetHisFileName(lastHisRecID)

		//保持历史记录到CSV文件
		err:=SaveRecsToCSV(hisFileName,hisRecItems)
		if err!=nil {
			log.Println("end PredictController Predict with error:SaveHisRecsToCSV error")
			rsp:=common.CreateResponse(common.CreateError(common.ResultSaveHisrecToCSVError,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		//调用预测服务
		resultFileName:=GetResultFileName(lastHisRecID)
		result:=Predict(hisFileName,resultFileName)
		if result!="0" {
			log.Println("end PredictController Predict with error:Predict error")
			params:=map[string]interface{}{
				"message":result,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultPredictError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		//读取结果数据
		resultRecs,err:=ReadResultRecsFromCSV(resultFileName)
		if err!=nil {
			log.Println("end PredictController Predict with error:ReadResultRecsFromCSV error")
			rsp:=common.CreateResponse(common.CreateError(common.ResultReadResultRecsFromCSVError,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

		//保存结果数据到数据库
		errorCode=SaveRecsToDB(resultRecs,pc.CRVClient,header.Token)
		if errorCode!=common.ResultSuccess {
			log.Println("end PredictController Predict with error:SavePredictToDB error")
			rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}

	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
}