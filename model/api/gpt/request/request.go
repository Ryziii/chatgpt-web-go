package request

type ChatProcessRequest struct {
	Prompt        string      `json:"prompt"`
	Options       ChatContext `json:"options,omitempty"`
	SystemMessage string      `json:"systemMessage"`
	Temperature   float64     `json:"temperature,omitempty"`
	TopP          float64     `json:"top_p,omitempty"`
}

type ChatContext struct {
	ConversationId  string `json:"conversationId,omitempty"`
	ParentMessageId string `json:"parentMessageId,omitempty"`
}
