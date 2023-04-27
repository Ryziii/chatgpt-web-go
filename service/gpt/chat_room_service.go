package gpt

import (
	"chatgpt-web-go/global"
	gpt2 "chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/repository"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
)

type ChatRoomService interface {
	CreateChatRoom(chatMessageDO *gpt2.ChatMessageDO) (gpt2.ChatRoomDO, error)
}

type chatRoomService struct {
	chatRoomRepo repository.ChatRoomRepository
}

func (s *chatRoomService) CreateChatRoom(chatMessageDO *gpt2.ChatMessageDO) (gpt2.ChatRoomDO, error) {
	chatRoom := gpt2.ChatRoomDO{
		ApiType:            chatMessageDO.APIType,
		IP:                 "",
		FirstChatMessageID: chatMessageDO.ID,
		ConversationID:     chatMessageDO.ConversationID,
		FirstMessageID:     uuid.New().String(),
		Title: func() string {
			ru := []rune(chatMessageDO.Content)
			return string(ru[:int(math.Min(float64(len(ru)), 50))])
		}(),
	}

	err := s.chatRoomRepo.CreateChatRoom(&chatRoom)
	if err != nil {
		global.Gzap.Error("CreateChatRoom", zap.Error(err))
		return gpt2.ChatRoomDO{}, errors.New("系统内部错误, 请联系管理员")
	}

	return chatRoom, nil
}

func NewChatRoomService() ChatRoomService {
	return &chatRoomService{chatRoomRepo: repository.NewChatRoomRepository()}
}
