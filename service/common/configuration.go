package common

import (
	"log"
	"os"
	"encoding/json"
)

type predictConf struct {
	HisCount int `json:"hisCount"`
}

type dataCleanConf struct {
	HisCount int `json:"hisCount"`
	FollowingCount int `json:"followingCount"`
}

type serviceConf struct {
	Port string `json:"port"`
}

type crvConf struct {
	Server string `json:"server"`
	Token string `json:"token"`
}

type Config struct {
	Service serviceConf `json:"service"`
	CRV crvConf `json:"crv"`
	Predict predictConf `json:"predict"`
	DataClean dataCleanConf `json:"dataClean"`
}

var gConfig Config

func InitConfig()(*Config){
	log.Println("init configuation start ...")
	//获取用户账号
	//获取用户角色信息
	//根据角色过滤出功能列表
	fileName := "conf/conf.json"
	filePtr, err := os.Open(fileName)
	if err != nil {
        log.Fatal("Open file failed [Err:%s]", err.Error())
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
	}
	log.Println("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}