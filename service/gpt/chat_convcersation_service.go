package gpt

import (
	"chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/repository"
)

type ChatConversationService interface {
	CreateConversation(chatConversation *gpt.ChatConversation) error
}

type chatConversationService struct {
	chatConversationRepo repository.ChatConversationRepository
}

func NewChatConversationService() ChatConversationService {
	return &chatConversationService{chatConversationRepo: repository.NewChatConversationRepository()}
}

func (s *chatConversationService) CreateConversation(chatConversation *gpt.ChatConversation) error {
	return s.chatConversationRepo.CreateChatConversation(chatConversation)
}
