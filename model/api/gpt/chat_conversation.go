package gpt

import (
	"chatgpt-web-go/model/common"
)

type ChatConversation struct {
	common.Model
	QuestionId uint64
	AnswerId   uint64
	ChatRoomId uint64
	ParentId   uint64 `gorm:"column:parent_conversation_id"`
}

func (ChatConversation) TableName() string {
	return "chat_conversation"
}
