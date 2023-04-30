package response

type ChatReplyMessage struct {
	Role string `json:"role"`

	Id string `json:"id"`

	ParentMessageId string `json:"parentMessageId"`

	ConversationId string `json:"conversationId"`

	Text string `json:"text"`
}
