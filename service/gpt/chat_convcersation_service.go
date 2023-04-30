package gpt

import (
	"chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/model/api/gpt/request"
	"chatgpt-web-go/model/common"
	"chatgpt-web-go/repository"
	"chatgpt-web-go/utils"
	"strconv"
)

type ChatConversationService interface {
	CreateConversation(chatConversation *gpt.ChatConversation) error
	InitChatConversation(*gpt.ChatConversation, request.ChatProcessRequest) error
	GetConversationById(id uint64, con *gpt.ChatConversation) error
}

type chatConversationService struct {
	chatConversationRepo repository.ChatConversationRepository
	chatRoomService      ChatRoomService
}

func (s *chatConversationService) GetConversationById(id uint64, con *gpt.ChatConversation) error {
	return s.chatConversationRepo.GetOne(con, gpt.ChatConversation{Model: common.Model{Id: id}})
}

func NewChatConversationService() ChatConversationService {
	return &chatConversationService{chatConversationRepo: repository.NewChatConversationRepository(), chatRoomService: NewChatRoomService()}
}
func (s *chatConversationService) InitChatConversation(conversation *gpt.ChatConversation, req request.ChatProcessRequest) error {
	// 一刀切，conversationId为空或转换失败，视为没有父对话
	conversationId, _ := strconv.ParseUint(req.Options.ParentMessageId, 10, 64)
	if conversationId == 0 {
		chatRoom, _ := s.chatRoomService.CreateChatRoom()
		conversation.ChatRoomId = chatRoom.Id
		conversation.ChatRoom = &chatRoom
	} else {
		if err := s.chatConversationRepo.GetChatConversationById(conversation, conversationId); err != nil {
			return err
		} else {
			conversation.ParentId = conversationId
		}
	}
	conversation.Id = utils.GetSnowIdUint64()
	conversation.QuestionId = utils.GetSnowIdUint64()
	conversation.AnswerId = utils.GetSnowIdUint64()
	conversation.Question = nil
	conversation.Answer = nil
	conversation.ContextCount++
	return nil
}
func (s *chatConversationService) CreateConversation(chatConversation *gpt.ChatConversation) error {
	chatConversation.ContextCount++
	chatConversation.QuestionUseToken = chatConversation.Question.TotalTokens
	chatConversation.AnswerUseToken = chatConversation.Answer.TotalTokens
	chatConversation.TotalTokens += chatConversation.AnswerUseToken + chatConversation.QuestionUseToken
	return s.chatConversationRepo.CreateChatConversation(chatConversation)
}
