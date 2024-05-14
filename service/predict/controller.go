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
}

func (pc *PredictController) Bind(router *gin.Engine) {
	router.POST("/predict", pc.Predict)
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

	lastHisRecID:=(*rep.SelectedRowKeys)[0]
	log.Println("lastHisRecID:",lastHisRecID)



	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
}