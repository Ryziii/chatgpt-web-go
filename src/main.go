package main

import (
	"chatgpt-web-go/src/global"
	"chatgpt-web-go/src/initialize"
	"chatgpt-web-go/src/routes"
	"net/http"
)

func main() {
	initialize.InitViper()
	initialize.InitGormMysql()
	initialize.InitZap()
	initialize.InitRedis()
	initialize.InitGPT()
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
