package initialize

import (
	"chatgpt-web-go/global"
	"github.com/sashabaranov/go-openai"
)

func InitGPT() {
	config := openai.DefaultConfig(global.Cfg.GPT.Token)
	config.BaseURL = global.Cfg.GPT.BaseURL
	global.GPTClient = openai.NewClientWithConfig(config)
}
