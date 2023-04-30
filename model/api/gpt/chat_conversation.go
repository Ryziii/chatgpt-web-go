package gpt

import (
	"chatgpt-web-go/model/common"
)

type ChatConversation struct {
	common.Model
	QuestionId uint64
	Question   ChatMessage `json:"-" gorm:"foreignKey:QuestionId;references:Id"`
	AnswerId   uint64
	Answer     ChatMessage `json:"-" gorm:"foreignKey:AnswerId;references:Id"`
	ChatRoomId uint64
	ChatRoom   ChatRoom `json:"-" gorm:"foreignKey:ChatRoomId;references:Id"`
	ParentId   uint64   `gorm:"column:parent_conversation_id"`
}

func (ChatConversation) TableName() string {
	return "chat_conversation"
}
