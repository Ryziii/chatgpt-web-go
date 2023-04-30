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
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ChatMessageService interface {
	SaveChatMessage(*model.ChatMessage) error
	GetOpenAiRequest(request.ChatProcessRequest, *model.ChatConversation) (openai.ChatCompletionRequest, error)
	UpdateChatMessage(answer *model.ChatMessage) error
}

type chatMessageService struct {
	TotalToken              int
	SystemMessage           openai.ChatCompletionMessage
	chatMessageRepo         repository.ChatMessageRepository
	chatConversationService ChatConversationService
}

func (s *chatMessageService) UpdateChatMessage(answer *model.ChatMessage) error {
	if err := s.chatMessageRepo.UpdateChatMessage(answer); err != nil {
		return err
	}
	return nil
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
	times := 0
	buildMessages = func(conv *model.ChatConversation) {
		times++
		if global.Cfg.GPT.RecurveTimes == times {
			return
		}
		if conv.Answer != nil && conv.Question != nil && (conv.Answer.Status == enum.PART_SUCCESS || conv.Answer.Status == enum.COMPLETE_SUCCESS) {
			if err := addMessages(conv); err != nil {
				return
			}
		}
		if conv.ParentId == 0 {
			return
		}
		var parConv model.ChatConversation
		if err := s.chatConversationService.GetConversationById(conv.ParentId, &parConv); err != nil {
			return
		}
		buildMessages(&parConv)
	}

	reCon := new(model.ChatConversation)
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
	conversation.Question = &model.ChatMessage{
		Model:     common.Model{Id: conversation.QuestionId},
		ModelName: completionRequest.Model,
		Content:   req.Prompt,
		OriginalData: func() string {
			jsonV, _ := json.Marshal(completionRequest)
			return string(jsonV)
		}(),
		TotalTokens: s.TotalToken,
		Status:      enum.COMPLETE_SUCCESS,
	}
	conversation.Answer = &model.ChatMessage{
		Model:     common.Model{Id: conversation.AnswerId},
		ModelName: completionRequest.Model,
		Status:    enum.INIT,
	}

	return completionRequest, nil
}

func (s *chatMessageService) SaveChatMessage(chatMessage *model.ChatMessage) error {
	if err := s.chatMessageRepo.CreateChatMessage(chatMessage); err != nil {
		return err
	}
	return nil
}

func NewChatMessageService() ChatMessageService {
	return &chatMessageService{
		chatMessageRepo:         repository.NewChatMessageRepository(),
		chatConversationService: NewChatConversationService(),
		SystemMessage:           getDefaultSystemMessage(),
		TotalToken:              utils.NumTokensFromMessages([]openai.ChatCompletionMessage{getDefaultSystemMessage()}, openai.GPT3Dot5Turbo) + global.Cfg.GPT.MaxToken,
	}
}

func (s *chatMessageService) updateTotalTokens(tokens int) bool {
	s.TotalToken += tokens
	if s.TotalToken > 4000 {
		return false
	}
	return true
}

func getDefaultSystemMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: "As an all-knowing assistant, answer any question concisely and accurately in Chinese"}
}
