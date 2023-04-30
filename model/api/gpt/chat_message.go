package gpt

import (
	enum "chatgpt-web-go/global/enum/gpt"
	"chatgpt-web-go/model/common"
)

type ChatMessage struct {
	common.Model
	ModelName         string
	IP                string
	Content           string
	OriginalData      string
	ResponseErrorData string
	TotalTokens       int
	Status            enum.ChatMessageStatusEnum
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
