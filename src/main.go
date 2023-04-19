package main

import (
	"chatgpt-web-go/src/global"
	"chatgpt-web-go/src/initialize"
	"chatgpt-web-go/src/routes"
	"net/http"
)

func main() {
	initialize.ViperInit()
	initialize.GormMysqlInit()
	initialize.InitZap()
	initialize.InitRedis()
	handler := routes.InitApiRoutes()
	server := &http.Server{
		Addr:    ":3002",
		Handler: handler,
	}
	if global.Gdb != nil {
		dbConfig, _ := global.Gdb.DB()
		defer dbConfig.Close()
	}
	if global.Gzap != nil {
		defer global.Gzap.Sync()
	}
	if global.Gredis != nil {
		defer global.Gredis.Close()
	}
	server.ListenAndServe()
}
