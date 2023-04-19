package initialize

import (
	"chatgpt-web-go/src/global"
	"go.uber.org/zap"
)

func InitZap() {
	global.Gzap, _ = zap.NewDevelopment()
}
