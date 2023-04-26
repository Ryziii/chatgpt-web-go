package gpt

import (
	"chatgpt-web-go/src/global"
	"chatgpt-web-go/src/model/api/gpt"
	"chatgpt-web-go/src/repository"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
)

type ChatRoomService interface {
	CreateChatRoom(chatMessageDO *gpt.ChatMessageDO) (gpt.ChatRoomDO, error)
}

type chatRoomService struct {
	chatRoomRepo repository.ChatRoomRepository
}

func (s *chatRoomService) CreateChatRoom(chatMessageDO *gpt.ChatMessageDO) (gpt.ChatRoomDO, error) {
	chatRoom := gpt.ChatRoomDO{
		ApiType:            chatMessageDO.APIType,
		IP:                 "",
		FirstChatMessageID: chatMessageDO.ID,
		ConversationID:     chatMessageDO.ConversationID,
		FirstMessageID:     uuid.New().String(),
		Title:              chatMessageDO.Content[:int(math.Min(float64(len(chatMessageDO.Content)), 50))],
	}

	err := s.chatRoomRepo.CreateChatRoom(&chatRoom)
	if err != nil {
		global.Gzap.Error("CreateChatRoom", zap.Error(err))
		return gpt.ChatRoomDO{}, errors.New("系统内部错误, 请联系管理员")
	}

	return chatRoom, nil
}

func NewChatRoomService() ChatRoomService {
	return &chatRoomService{chatRoomRepo: repository.NewChatRoomRepository()}
}
