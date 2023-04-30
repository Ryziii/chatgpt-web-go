package gpt

import (
	gpt2 "chatgpt-web-go/global/enum/gpt"
	"chatgpt-web-go/model/common"
	"fmt"
)

type ChatMessage struct {
	common.Model
	MessageId               string
	ParentMessageId         string
	ParentAnswerMessageId   string
	ParentQuestionMessageId string
	ContextCount            int
	QuestionContextCount    int
	MessageType             gpt2.ChatMessageTypeEnum
	ChatRoomId              uint64
	ConversationId          string
	APIType                 gpt2.ApiTypeEnum
	ModelName               string
	IP                      string
	APIKey                  string
	Content                 string
	OriginalData            string
	ResponseErrorData       string
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	Status                  gpt2.ChatMessageStatusEnum
	IsHide                  bool
}

func (ChatMessage) TableName() string {
	return "chat_message"
}

func (c ChatMessage) ToString() string {
	return fmt.Sprintf("Id: %d\n "+
		"MessageId: %s\n "+
		"ParentMessageId: %s\n "+
		"ParentAnswerMessageId: %s\n "+
		"ParentQuestionMessageId: %s\n "+
		"ContextCount: %d\n "+
		"QuestionContextCount: %d\n "+
		"MessageType: %s\n "+
		"ChatRoomId: %d\n "+
		"ConversationId: %s\n "+
		"APIType: %s\n "+
		"ModelName: %s\n "+
		"IP: %s\n "+
		"APIKey: %s\n "+
		"Content: %s\n "+
		"OriginalData: %s\n "+
		"ResponseErrorData: %s\n "+
		"PromptTokens: %d\n "+
		"CompletionTokens: %d\n "+
		"TotalTokens: %d\n "+
		"Status: %s\n "+
		"IsHide: %t\n "+
		"CreatedTime: %s\n "+
		"UpdatedTime: %s\n", c.Id, c.MessageId, c.ParentMessageId, c.ParentAnswerMessageId, c.ParentQuestionMessageId, c.ContextCount, c.QuestionContextCount, c.MessageType, c.ChatRoomId, c.ConversationId, c.APIType, c.ModelName, c.IP, c.APIKey, c.Content, c.OriginalData, c.ResponseErrorData, c.PromptTokens, c.CompletionTokens, c.TotalTokens, c.Status, c.IsHide, c.CreateTime, c.UpdateTime)
}
