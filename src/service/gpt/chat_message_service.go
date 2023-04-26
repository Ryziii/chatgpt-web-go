package gpt

import (
	"chatgpt-web-go/src/global"
	enum "chatgpt-web-go/src/global/enum/gpt"
	"chatgpt-web-go/src/model/api/gpt"
	"chatgpt-web-go/src/model/api/gpt/request"
	"chatgpt-web-go/src/repository"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ChatMessageService interface {
	InitChatMessage(chatMessageDO *gpt.ChatMessageDO, chatProcessRequest request.ChatProcessRequest, apiTypeEnum enum.ApiTypeEnum) error
	PopulateInitParentMessage(chatMessageDO *gpt.ChatMessageDO, chatProcessRequest request.ChatProcessRequest) error
	addContextChatMessage(chatMessageDO *gpt.ChatMessageDO, messages *[]openai.ChatCompletionMessage)
	GetOpenAiRequestReady(req request.ChatProcessRequest) (gpt.ChatMessageDO, openai.ChatCompletionRequest, error)
}

type chatMessageService struct {
	chatMessageRepo repository.ChatMessageRepository
}

func NewChatMessageService() ChatMessageService {
	return &chatMessageService{chatMessageRepo: repository.NewChatMessageRepository()}
}

func (s *chatMessageService) InitChatMessage(chatMessageDO *gpt.ChatMessageDO, chatProcessRequest request.ChatProcessRequest, apiTypeEnum enum.ApiTypeEnum) error {

	*chatMessageDO = gpt.ChatMessageDO{
		Model: gpt.Model{ID: uint64(func() int64 {
			snowNode, _ := snowflake.NewNode(1)
			id := snowNode.Generate().Int64()
			return id
		}())},
		MessageID:        uuid.New().String(),
		ConversationID:   uuid.New().String(),
		MessageType:      enum.QUESTION,
		APIType:          apiTypeEnum,
		Content:          chatProcessRequest.Prompt,
		ModelName:        global.Cfg.GPT.OpenAIAPIMODEL,
		OriginalData:     "",
		PromptTokens:     -1,
		CompletionTokens: -1,
		TotalTokens:      -1,
		IP:               "",
		Status:           enum.INIT,
	}

	if err := s.PopulateInitParentMessage(chatMessageDO, chatProcessRequest); err != nil {
		return err
	}

	return nil
}

func (s *chatMessageService) PopulateInitParentMessage(chatMessageDO *gpt.ChatMessageDO, chatProcessRequest request.ChatProcessRequest) error {
	parentMessageID := chatProcessRequest.Options.ParentMessageID
	conversationID := chatProcessRequest.Options.ConversationID

	if parentMessageID != "" && conversationID != "" {
		parentChatMessage := gpt.ChatMessageDO{}
		err := s.chatMessageRepo.GetOne(&parentChatMessage, gpt.ChatMessageDO{
			MessageID:      parentMessageID,
			ConversationID: conversationID,
			APIType:        chatMessageDO.APIType,
			MessageType:    enum.ANSWER,
		})
		chatMessageDO.ParentMessageID = parentMessageID
		chatMessageDO.ConversationID = conversationID
		chatMessageDO.ParentAnswerMessageID = parentMessageID
		if err != nil || parentChatMessage == (gpt.ChatMessageDO{}) {
			return errors.New("系统出错, 无法找到聊天记录. 请尝试关闭输入框左边的携带聊天记录按钮后重试, 或新建聊天.")
		}
		chatMessageDO.ParentQuestionMessageID = parentChatMessage.ParentQuestionMessageID
		chatMessageDO.ChatRoomID = parentChatMessage.ChatRoomID
		chatMessageDO.ContextCount = parentChatMessage.ContextCount + 1
		chatMessageDO.QuestionContextCount = parentChatMessage.QuestionContextCount + 1

		if chatMessageDO.APIType == enum.AccessToken {
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

func (s *chatMessageService) addContextChatMessage(chatMessageDO *gpt.ChatMessageDO, messages *[]openai.ChatCompletionMessage) {
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
	if chatMessageDO.MessageType == enum.ANSWER {
		role = openai.ChatMessageRoleAssistant
	} else {
		role = openai.ChatMessageRoleUser
	}

	// 答案不是成功的, 不加入messages
	if chatMessageDO.MessageType == enum.ANSWER && (chatMessageDO.Status != enum.PART_SUCCESS && chatMessageDO.Status != enum.COMPLETE_SUCCESS) {
		if chatMessageDO.ParentAnswerMessageID == "" {
			return
		}
		parentMessage := gpt.ChatMessageDO{}
		err := s.chatMessageRepo.GetOne(&parentMessage, gpt.ChatMessageDO{
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
	parentMessage := gpt.ChatMessageDO{}
	err := s.chatMessageRepo.GetOne(&parentMessage, gpt.ChatMessageDO{
		MessageID: chatMessageDO.ParentMessageID,
	})
	if err != nil {
		return
	}
	s.addContextChatMessage(&parentMessage, messages)
	return
}

func (s *chatMessageService) GetOpenAiRequestReady(req request.ChatProcessRequest) (gpt.ChatMessageDO, openai.ChatCompletionRequest, error) {
	var chatMessageDO gpt.ChatMessageDO
	var completionRequest openai.ChatCompletionRequest
	global.Gzap.Info("chatMessageDO.ID", zap.Any("chatMessageDO.ID", chatMessageDO.ID))
	if err := s.InitChatMessage(&chatMessageDO, req, enum.ApiKey); err != nil {
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
