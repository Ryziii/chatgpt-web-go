package main

import (
	"chatgpt-web-go/global"
	initialize2 "chatgpt-web-go/initialize"
	"chatgpt-web-go/routes"
	"net/http"
)

func main() {
	initialize2.InitViper()
	initialize2.InitGormMysql()
	initialize2.InitZap()
	initialize2.InitRedis()
	initialize2.InitGPT()
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
