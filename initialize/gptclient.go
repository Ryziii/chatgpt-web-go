package initialize

import (
	"chatgpt-web-go/global"
	"github.com/sashabaranov/go-openai"
)

func InitGPT() {
	global.GPTClient = openai.NewClient(global.Cfg.GPT.Token)
}
