package initialize

import (
	"chatgpt-web-go/src/global"
	"github.com/sashabaranov/go-openai"
)

func InitGPT() {
	global.GPTClient = openai.NewClient(global.Cfg.GPT.Token)
}
