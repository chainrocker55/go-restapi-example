package main

import (
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/database"
	"creditlimit-connector/app/log"
	"creditlimit-connector/app/server"
	"os"
	_ "time/tzdata"
)

func main() {
	os.Setenv("TZ", "Asia/Bangkok")
	configs.InitConfig()
	log.Init()
	database.InitINVXCreditLimitDatabase()
	database.InitINVXCenterDatabase()
	database.InitRedisClient()
	server.Serve()
}
