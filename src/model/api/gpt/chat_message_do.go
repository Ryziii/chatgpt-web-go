package gpt

import (
	"chatgpt-web-go/src/global/enum/gpt"
	"fmt"
)

type ChatMessageDO struct {
	Model
	MessageID               string
	ParentMessageID         string
	ParentAnswerMessageID   string
	ParentQuestionMessageID string
	ContextCount            int
	QuestionContextCount    int
	MessageType             gpt.ChatMessageTypeEnum
	ChatRoomID              uint64
	ConversationID          string
	APIType                 gpt.ApiTypeEnum
	ModelName               string
	IP                      string
	APIKey                  string
	Content                 string
	OriginalData            string
	ResponseErrorData       string
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	Status                  gpt.ChatMessageStatusEnum
	IsHide                  bool
}

func (ChatMessageDO) TableName() string {
	return "chat_message"
}

func (c ChatMessageDO) ToString() string {
	return fmt.Sprintf("ID: %d\n "+
		"MessageID: %s\n "+
		"ParentMessageID: %s\n "+
		"ParentAnswerMessageID: %s\n "+
		"ParentQuestionMessageID: %s\n "+
		"ContextCount: %d\n "+
		"QuestionContextCount: %d\n "+
		"MessageType: %s\n "+
		"ChatRoomID: %d\n "+
		"ConversationID: %s\n "+
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
		"UpdatedTime: %s\n", c.ID, c.MessageID, c.ParentMessageID, c.ParentAnswerMessageID, c.ParentQuestionMessageID, c.ContextCount, c.QuestionContextCount, c.MessageType, c.ChatRoomID, c.ConversationID, c.APIType, c.ModelName, c.IP, c.APIKey, c.Content, c.OriginalData, c.ResponseErrorData, c.PromptTokens, c.CompletionTokens, c.TotalTokens, c.Status, c.IsHide, c.CreateTime, c.UpdateTime)
}
