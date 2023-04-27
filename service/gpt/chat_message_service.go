package gpt

import (
	"chatgpt-web-go/global"
	gpt2 "chatgpt-web-go/global/enum/gpt"
	gpt3 "chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/model/api/gpt/request"
	"chatgpt-web-go/repository"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type ChatMessageService interface {
	InitChatMessage(chatMessageDO *gpt3.ChatMessageDO, chatProcessRequest request.ChatProcessRequest, apiTypeEnum gpt2.ApiTypeEnum) error
	PopulateInitParentMessage(chatMessageDO *gpt3.ChatMessageDO, chatProcessRequest request.ChatProcessRequest) error
	addContextChatMessage(chatMessageDO *gpt3.ChatMessageDO, messages *[]openai.ChatCompletionMessage)
	GetOpenAiRequestReady(req request.ChatProcessRequest) (gpt3.ChatMessageDO, openai.ChatCompletionRequest, error)
}

type chatMessageService struct {
	chatMessageRepo repository.ChatMessageRepository
}

func NewChatMessageService() ChatMessageService {
	return &chatMessageService{chatMessageRepo: repository.NewChatMessageRepository()}
}

func (s *chatMessageService) InitChatMessage(chatMessageDO *gpt3.ChatMessageDO, chatProcessRequest request.ChatProcessRequest, apiTypeEnum gpt2.ApiTypeEnum) error {

	*chatMessageDO = gpt3.ChatMessageDO{
		Model: gpt3.Model{ID: uint64(func() int64 {
			snowNode, _ := snowflake.NewNode(1)
			id := snowNode.Generate().Int64()
			return id
		}())},
		MessageID:        uuid.New().String(),
		ConversationID:   uuid.New().String(),
		MessageType:      gpt2.QUESTION,
		APIType:          apiTypeEnum,
		Content:          chatProcessRequest.Prompt,
		ModelName:        global.Cfg.GPT.OpenAIAPIMODEL,
		OriginalData:     "",
		PromptTokens:     -1,
		CompletionTokens: -1,
		TotalTokens:      -1,
		IP:               "",
		Status:           gpt2.INIT,
	}

	if err := s.PopulateInitParentMessage(chatMessageDO, chatProcessRequest); err != nil {
		return err
	}

	return nil
}

func (s *chatMessageService) PopulateInitParentMessage(chatMessageDO *gpt3.ChatMessageDO, chatProcessRequest request.ChatProcessRequest) error {
	parentMessageID := chatProcessRequest.Options.ParentMessageID
	conversationID := chatProcessRequest.Options.ConversationID

	if parentMessageID != "" && conversationID != "" {
		parentChatMessage := gpt3.ChatMessageDO{}
		err := s.chatMessageRepo.GetOne(&parentChatMessage, gpt3.ChatMessageDO{
			MessageID:      parentMessageID,
			ConversationID: conversationID,
			APIType:        chatMessageDO.APIType,
			MessageType:    gpt2.ANSWER,
		})
		chatMessageDO.ParentMessageID = parentMessageID
		chatMessageDO.ConversationID = conversationID
		chatMessageDO.ParentAnswerMessageID = parentMessageID
		if err != nil || parentChatMessage == (gpt3.ChatMessageDO{}) {
			return errors.New("系统出错, 无法找到聊天记录. 请尝试关闭输入框左边的携带聊天记录按钮后重试, 或新建聊天.")
		}
		chatMessageDO.ParentQuestionMessageID = parentChatMessage.ParentQuestionMessageID
		chatMessageDO.ChatRoomID = parentChatMessage.ChatRoomID
		chatMessageDO.ContextCount = parentChatMessage.ContextCount + 1
		chatMessageDO.QuestionContextCount = parentChatMessage.QuestionContextCount + 1

		if chatMessageDO.APIType == gpt2.AccessToken {
			if chatMessageDO.ModelName != parentChatMessage.ModelName {
				return errors.New("model name not consistent with parent message")
			}
		}

	} else {
		chatr := NewChatRoomService()
		chatRoomDO, err := chatr.CreateChatRoom(chatMessageDO)
		if err != nil {
			return err
		}
		chatMessageDO.ChatRoomID = chatRoomDO.ID
		chatMessageDO.ContextCount = 1
		chatMessageDO.QuestionContextCount = 1
	}
	return nil
}

func (s *chatMessageService) addContextChatMessage(chatMessageDO *gpt3.ChatMessageDO, messages *[]openai.ChatCompletionMessage) {
	if chatMessageDO == nil {
		return
	}
	if chatMessageDO.ParentMessageID == "" {
		*messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: chatMessageDO.Content,
			},
		}, *messages...)
		return
	}

	var role string
	if chatMessageDO.MessageType == gpt2.ANSWER {
		role = openai.ChatMessageRoleAssistant
	} else {
		role = openai.ChatMessageRoleUser
	}

	// 答案不是成功的, 不加入messages
	if chatMessageDO.MessageType == gpt2.ANSWER && (chatMessageDO.Status != gpt2.PART_SUCCESS && chatMessageDO.Status != gpt2.COMPLETE_SUCCESS) {
		if chatMessageDO.ParentAnswerMessageID == "" {
			return
		}
		parentMessage := gpt3.ChatMessageDO{}
		err := s.chatMessageRepo.GetOne(&parentMessage, gpt3.ChatMessageDO{
			MessageID: chatMessageDO.ParentAnswerMessageID,
		})
		if err != nil {
			return
		}
		s.addContextChatMessage(&parentMessage, messages)
		return
	}
	*messages = append([]openai.ChatCompletionMessage{
		{
			Role:    role,
			Content: chatMessageDO.Content,
		},
	}, *messages...)
	parentMessage := gpt3.ChatMessageDO{}
	err := s.chatMessageRepo.GetOne(&parentMessage, gpt3.ChatMessageDO{
		MessageID: chatMessageDO.ParentMessageID,
	})
	if err != nil {
		return
	}
	s.addContextChatMessage(&parentMessage, messages)
	return
}

func (s *chatMessageService) GetOpenAiRequestReady(req request.ChatProcessRequest) (gpt3.ChatMessageDO, openai.ChatCompletionRequest, error) {
	var chatMessageDO gpt3.ChatMessageDO
	var completionRequest openai.ChatCompletionRequest
	if err := s.InitChatMessage(&chatMessageDO, req, gpt2.ApiKey); err != nil {
		return chatMessageDO, completionRequest, err
	}

	var messages []openai.ChatCompletionMessage
	s.addContextChatMessage(&chatMessageDO, &messages)

	if req.SystemMessage != "" {
		systemMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemMessage,
		}
		messages = append([]openai.ChatCompletionMessage{systemMessage}, messages...)
	}

	completionRequest = openai.ChatCompletionRequest{
		Model:           global.Cfg.GPT.OpenAIAPIMODEL,
		Messages:        messages,
		MaxTokens:       global.Cfg.GPT.MaxToken,
		Temperature:     global.Cfg.GPT.Temperature,
		TopP:            global.Cfg.GPT.TopP,
		N:               1,
		Stream:          true,
		PresencePenalty: 1,
	}
	return chatMessageDO, completionRequest, nil
}
