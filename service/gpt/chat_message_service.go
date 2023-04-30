package gpt

import (
	"chatgpt-web-go/global"
	enum "chatgpt-web-go/global/enum/gpt"
	model "chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/model/api/gpt/request"
	"chatgpt-web-go/model/common"
	"chatgpt-web-go/repository"
	"chatgpt-web-go/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ChatMessageService interface {
	GetOpenAiRequestReady(req request.ChatProcessRequest) (model.ChatMessage, openai.ChatCompletionRequest, error)
	SaveQuestionDOFromChatMessage(string, model.ChatMessage, openai.ChatCompletionRequest) error
	GetOpenAiRequest(req request.ChatProcessRequest, conversation *model.ChatConversation) (openai.ChatCompletionRequest, error)
}

type chatMessageService struct {
	TotalToken              int
	SystemMessage           openai.ChatCompletionMessage
	chatMessageRepo         repository.ChatMessageRepository
	chatConversationService ChatConversationService
}

func (s *chatMessageService) GetOpenAiRequest(req request.ChatProcessRequest, conversation *model.ChatConversation) (openai.ChatCompletionRequest, error) {
	var completionRequest openai.ChatCompletionRequest
	var messages []openai.ChatCompletionMessage
	thisReqMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Prompt,
	}
	messages = append(messages, thisReqMessage)

	// 添加问答进message并校验是否超token
	addMessages := func(conv *model.ChatConversation) error {
		var rem []openai.ChatCompletionMessage
		question := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: conv.Question.Content,
		}
		answer := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: conv.Answer.Content,
		}
		rem = append(rem, answer)
		if !s.updateTotalTokens(utils.NumTokensFromMessages(rem, openai.GPT3Dot5Turbo)) {
			return errors.New("totalToken不足")
		}
		messages = append([]openai.ChatCompletionMessage{answer}, messages...)
		if !s.updateTotalTokens(utils.NumTokensFromMessages(rem, openai.GPT3Dot5Turbo)) {
			return errors.New("totalToken不足")
		}
		messages = append([]openai.ChatCompletionMessage{question}, messages...)
		return nil
	}
	// 递归conversation将上文填充进messages
	var buildMessages func(conv *model.ChatConversation)
	buildMessages = func(conv *model.ChatConversation) {
		if conv.Answer != nil && conv.Question != nil && (conv.Answer.Status == enum.PART_SUCCESS || conv.Answer.Status == enum.COMPLETE_SUCCESS) {
			if err := addMessages(conv); err != nil {
				return
			}
		}
		if conv.ParentId == 0 {
			return
		}
		if err := s.chatConversationService.GetConversationById(conv.ParentId, conv); err != nil {
			return
		}
		buildMessages(conv)
	}

	var reCon *model.ChatConversation
	if err := utils.DeepCopy(conversation, reCon); err != nil {
		global.Gzap.Error("buildMessages前深拷贝conversation错误", zap.Error(err))
		return openai.ChatCompletionRequest{}, err
	}
	buildMessages(reCon)

	messages = append([]openai.ChatCompletionMessage{s.SystemMessage}, messages...)
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
		return completionRequest, nil
	}

	return completionRequest, nil
}

func NewChatMessageService() ChatMessageService {
	return &chatMessageService{
		chatMessageRepo:         repository.NewChatMessageRepository(),
		chatConversationService: NewChatConversationService(),
		SystemMessage:           getDefaultSystemMessage(),
		TotalToken:              utils.NumTokensFromMessages([]openai.ChatCompletionMessage{getDefaultSystemMessage()}, openai.GPT3Dot5Turbo) + global.Cfg.GPT.MaxToken,
	}
}

func (s *chatMessageService) SaveQuestionDOFromChatMessage(ip string, chatMessageDO model.ChatMessage, completionRequest openai.ChatCompletionRequest) error {
	var questionDO model.ChatMessage
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
	questionDO.ParentAnswerMessageId = questionDO.ParentMessageId
	if err := s.chatMessageRepo.CreateChatMessage(&questionDO); err != nil {
		return err
	}
	return nil
}

func (s *chatMessageService) initChatMessage(chatMessageDO *model.ChatMessage, chatProcessRequest request.ChatProcessRequest, apiTypeEnum enum.ApiTypeEnum) error {
	*chatMessageDO = model.ChatMessage{
		Model:            common.Model{Id: utils.GetSnowIdUint64()},
		MessageId:        uuid.New().String(),
		ConversationId:   uuid.New().String(),
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

func (s *chatMessageService) populateInitParentMessage(chatMessageDO *model.ChatMessage, chatProcessRequest request.ChatProcessRequest) error {
	parentMessageId := chatProcessRequest.Options.ParentMessageId
	conversationId := chatProcessRequest.Options.ConversationId

	if parentMessageId != "" && conversationId != "" {
		parentChatMessage := model.ChatMessage{}
		err := s.chatMessageRepo.GetOne(&parentChatMessage, model.ChatMessage{
			MessageId:      parentMessageId,
			ConversationId: conversationId,
			APIType:        chatMessageDO.APIType,
			MessageType:    enum.ANSWER,
		})
		if err != nil || parentChatMessage == (model.ChatMessage{}) {
			return errors.New("系统出错, 无法找到聊天记录. 请尝试关闭输入框左边的携带聊天记录按钮后重试, 或新建聊天.")
		}
		chatMessageDO.ParentMessageId = parentMessageId
		chatMessageDO.ConversationId = conversationId
		chatMessageDO.ParentAnswerMessageId = parentMessageId
		chatMessageDO.ParentQuestionMessageId = parentChatMessage.ParentQuestionMessageId
		chatMessageDO.ChatRoomId = parentChatMessage.ChatRoomId
		chatMessageDO.ContextCount = parentChatMessage.ContextCount + 1
		chatMessageDO.QuestionContextCount = parentChatMessage.QuestionContextCount + 1
	} else {
		chatr := NewChatRoomService()
		chatRoomDO, err := chatr.CreateChatRoomByChatMessage(chatMessageDO)
		if err != nil {
			return err
		}
		chatMessageDO.ChatRoomId = chatRoomDO.Id
		chatMessageDO.ContextCount = 1
		chatMessageDO.QuestionContextCount = 1
	}
	return nil
}

func (s *chatMessageService) addContextChatMessage(chatMessageDO *model.ChatMessage, messages *[]openai.ChatCompletionMessage) {
	if chatMessageDO == nil {
		return
	}

	var processMessage func(chatMessageDO *model.ChatMessage)
	processMessage = func(chatMessageDO *model.ChatMessage) {
		// 没有父消息, 说明是第一条消息, 直接返回
		if chatMessageDO.ParentMessageId == "" {
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
			if chatMessageDO.ParentAnswerMessageId == "" {
				return
			}
			parentMessage := model.ChatMessage{}
			if err := s.chatMessageRepo.GetOne(&parentMessage, model.ChatMessage{
				MessageId: chatMessageDO.ParentAnswerMessageId,
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
		parentMessage := model.ChatMessage{}
		if err := s.chatMessageRepo.GetOne(&parentMessage, model.ChatMessage{
			MessageId: chatMessageDO.ParentMessageId,
		}); err != nil {
			return
		}
		processMessage(&parentMessage)
		return
	}

	processMessage(chatMessageDO)
}

func (s *chatMessageService) GetOpenAiRequestReady(req request.ChatProcessRequest) (model.ChatMessage, openai.ChatCompletionRequest, error) {
	//s.TotalToken = 0
	var chatMessageDO model.ChatMessage
	var completionRequest openai.ChatCompletionRequest
	if err := s.initChatMessage(&chatMessageDO, req, enum.ApiKey); err != nil {
		return chatMessageDO, completionRequest, err
	}

	var messages []openai.ChatCompletionMessage
	systemMessage := s.SystemMessage
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

func getDefaultSystemMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: "As an all-knowing and omnipotent assistant, you understand everything and can answer any question or solve any problem. Please provide concise, accurate, and brief answers in Chinese."}
}
