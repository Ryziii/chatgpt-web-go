package response

type ChatReplyMessage struct {
	Role string `json:"role"`

	ID string `json:"id"`

	ParentMessageID string `json:"parentMessageId"`

	ConversationID string `json:"conversationId"`

	Text string `json:"text"`
}
