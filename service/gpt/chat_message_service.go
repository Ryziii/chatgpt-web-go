package gpt

import (
	"chatgpt-web-go/global"
	enum "chatgpt-web-go/global/enum/gpt"
	model "chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/model/api/gpt/request"
	"chatgpt-web-go/repository"
	"chatgpt-web-go/utils"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type ChatMessageService interface {
	GetOpenAiRequestReady(req request.ChatProcessRequest) (model.ChatMessageDO, openai.ChatCompletionRequest, error)
	SaveQuestionDOFromChatMessageDO(string, model.ChatMessageDO, openai.ChatCompletionRequest) error
}

type chatMessageService struct {
	TotalToken      int
	chatMessageRepo repository.ChatMessageRepository
}

func NewChatMessageService() ChatMessageService {
	return &chatMessageService{
		chatMessageRepo: repository.NewChatMessageRepository(),
		TotalToken:      global.Cfg.GPT.MaxToken,
	}
}

func (s *chatMessageService) SaveQuestionDOFromChatMessageDO(ip string, chatMessageDO model.ChatMessageDO, completionRequest openai.ChatCompletionRequest) error {
	var questionDO model.ChatMessageDO
	if err := utils.DeepCopyByJson(&chatMessageDO, &questionDO); err != nil {
		return err
	}

	questionDO.IP = ip
	questionDO.OriginalData = func() string {
		jsonV, _ := json.Marshal(completionRequest)
		return string(jsonV)
	}()
	questionDO.PromptTokens = s.TotalToken
	questionDO.Status = enum.PART_SUCCESS
	questionDO.MessageType = enum.QUESTION
	questionDO.ParentAnswerMessageID = questionDO.ParentMessageID
	if err := s.chatMessageRepo.CreateChatMessage(&questionDO); err != nil {
		return err
	}
	return nil
}

func (s *chatMessageService) initChatMessageDO(chatMessageDO *model.ChatMessageDO, chatProcessRequest request.ChatProcessRequest, apiTypeEnum enum.ApiTypeEnum) error {
	*chatMessageDO = model.ChatMessageDO{
		Model: model.Model{ID: uint64(func() int64 {
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

	if err := s.populateInitParentMessage(chatMessageDO, chatProcessRequest); err != nil {
		return err
	}

	return nil
}

func (s *chatMessageService) populateInitParentMessage(chatMessageDO *model.ChatMessageDO, chatProcessRequest request.ChatProcessRequest) error {
	parentMessageID := chatProcessRequest.Options.ParentMessageID
	conversationID := chatProcessRequest.Options.ConversationID

	if parentMessageID != "" && conversationID != "" {
		parentChatMessage := model.ChatMessageDO{}
		err := s.chatMessageRepo.GetOne(&parentChatMessage, model.ChatMessageDO{
			MessageID:      parentMessageID,
			ConversationID: conversationID,
			APIType:        chatMessageDO.APIType,
			MessageType:    enum.ANSWER,
		})
		if err != nil || parentChatMessage == (model.ChatMessageDO{}) {
			return errors.New("系统出错, 无法找到聊天记录. 请尝试关闭输入框左边的携带聊天记录按钮后重试, 或新建聊天.")
		}
		chatMessageDO.ParentMessageID = parentMessageID
		chatMessageDO.ConversationID = conversationID
		chatMessageDO.ParentAnswerMessageID = parentMessageID
		chatMessageDO.ParentQuestionMessageID = parentChatMessage.ParentQuestionMessageID
		chatMessageDO.ChatRoomID = parentChatMessage.ChatRoomID
		chatMessageDO.ContextCount = parentChatMessage.ContextCount + 1
		chatMessageDO.QuestionContextCount = parentChatMessage.QuestionContextCount + 1
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

func (s *chatMessageService) addContextChatMessage(chatMessageDO *model.ChatMessageDO, messages *[]openai.ChatCompletionMessage) {
	if chatMessageDO == nil {
		return
	}

	var processMessage func(chatMessageDO *model.ChatMessageDO)
	processMessage = func(chatMessageDO *model.ChatMessageDO) {
		// 没有父消息, 说明是第一条消息, 直接返回
		if chatMessageDO.ParentMessageID == "" {
			message := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: chatMessageDO.Content,
			}
			token := utils.NumTokensFromMessages([]openai.ChatCompletionMessage{message}, openai.GPT3Dot5Turbo)
			if !s.updateTotalTokens(token) {
				return
			}
			*messages = append([]openai.ChatCompletionMessage{message}, *messages...)
			return
		}
		// 如果是回答, 但是状态不是成功, 寻找上一个回答
		if chatMessageDO.MessageType == enum.ANSWER && (chatMessageDO.Status != enum.PART_SUCCESS && chatMessageDO.Status != enum.COMPLETE_SUCCESS) {
			if chatMessageDO.ParentAnswerMessageID == "" {
				return
			}
			parentMessage := model.ChatMessageDO{}
			if err := s.chatMessageRepo.GetOne(&parentMessage, model.ChatMessageDO{
				MessageID: chatMessageDO.ParentAnswerMessageID,
			}); err != nil {
				return
			}
			processMessage(&parentMessage)
			return
		}
		message := openai.ChatCompletionMessage{
			Role: func() string {
				if chatMessageDO.MessageType == enum.ANSWER {
					return openai.ChatMessageRoleAssistant
				} else {
					return openai.ChatMessageRoleUser
				}
			}(),
			Content: chatMessageDO.Content,
		}
		token := utils.NumTokensFromMessages([]openai.ChatCompletionMessage{message}, openai.GPT3Dot5Turbo)
		if !s.updateTotalTokens(token) {
			return
		}
		*messages = append([]openai.ChatCompletionMessage{message}, *messages...)
		parentMessage := model.ChatMessageDO{}
		if err := s.chatMessageRepo.GetOne(&parentMessage, model.ChatMessageDO{
			MessageID: chatMessageDO.ParentMessageID,
		}); err != nil {
			return
		}
		processMessage(&parentMessage)
		return
	}

	processMessage(chatMessageDO)
}

func (s *chatMessageService) GetOpenAiRequestReady(req request.ChatProcessRequest) (model.ChatMessageDO, openai.ChatCompletionRequest, error) {
	//s.TotalToken = 0
	var chatMessageDO model.ChatMessageDO
	var completionRequest openai.ChatCompletionRequest
	if err := s.initChatMessageDO(&chatMessageDO, req, enum.ApiKey); err != nil {
		return chatMessageDO, completionRequest, err
	}

	var messages []openai.ChatCompletionMessage
	systemMessage := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: "As an all-knowing and omnipotent assistant, you understand everything and can answer any question or solve any problem. Please provide concise, accurate, and brief answers in Chinese."}
	token := utils.NumTokensFromMessages([]openai.ChatCompletionMessage{systemMessage}, openai.GPT3Dot5Turbo)
	s.updateTotalTokens(token)

	s.addContextChatMessage(&chatMessageDO, &messages)
	messages = append([]openai.ChatCompletionMessage{systemMessage}, messages...)

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
	numTokens := utils.NumTokensFromMessages(completionRequest.Messages, openai.GPT3Dot5Turbo)
	if ok := s.updateTotalTokens(numTokens); !ok {
		return chatMessageDO, completionRequest, nil
	}

	return chatMessageDO, completionRequest, nil
}

func (s *chatMessageService) updateTotalTokens(tokens int) bool {
	s.TotalToken += tokens
	if s.TotalToken > 4000 {
		return false
	}
	return true
}
