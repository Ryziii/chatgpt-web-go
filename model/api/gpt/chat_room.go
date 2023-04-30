package gpt

import (
	enum "chatgpt-web-go/global/enum/gpt"
	"chatgpt-web-go/model/common"
)

type ChatRoom struct {
	common.Model
	ConversationId     string
	IP                 string
	FirstChatMessageId uint64
	FirstMessageId     string
	Title              string
	ApiType            enum.ApiTypeEnum
}

func (ChatRoom) TableName() string {
	return "chat_room"
}
