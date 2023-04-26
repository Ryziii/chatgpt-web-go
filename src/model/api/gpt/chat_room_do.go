package gpt

import (
	enum "chatgpt-web-go/src/global/enum/gpt"
)

type ChatRoomDO struct {
	Model
	ConversationID     string
	IP                 string
	FirstChatMessageID uint64
	FirstMessageID     string
	Title              string
	ApiType            enum.ApiTypeEnum
}

func (ChatRoomDO) TableName() string {
	return "chat_room"
}
