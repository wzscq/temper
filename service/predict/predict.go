package predict

import (
	"os/exec"
	"log"
)

func Predict(hisFileName,resultFileName string)(string) {
	cmd := exec.Command("python3", "predict.py",hisFileName,resultFileName)
  	// 设置执行命令时的工作目录
  	_, err := cmd.Output()
	if err != nil {
		log.Println("Predict exec predict.py error:",err)
		return ""
	}
	//log.Println(string(out))
  	return "0"
}

func GetHisFileName(lastHisRecID string)(string){
	return lastHisRecID+"_his.csv"
}

func GetResultFileName(lastHisRecID string)(string){
	return lastHisRecID+"_result.csv"
}