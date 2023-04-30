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
	CreateChatRoomByChatMessage(chatMessageDO *gpt2.ChatMessage) (gpt2.ChatRoom, error)
	CreateChatRoom() (gpt2.ChatRoom, error)
}

type chatRoomService struct {
	chatRoomRepo repository.ChatRoomRepository
}

func (s *chatRoomService) CreateChatRoomByChatMessage(chatMessageDO *gpt2.ChatMessage) (gpt2.ChatRoom, error) {
	chatRoom := gpt2.ChatRoom{
		ApiType:            chatMessageDO.APIType,
		IP:                 "",
		FirstChatMessageId: chatMessageDO.Id,
		ConversationId:     chatMessageDO.ConversationId,
		FirstMessageId:     uuid.New().String(),
		Title: func() string {
			ru := []rune(chatMessageDO.Content)
			return string(ru[:int(math.Min(float64(len(ru)), 50))])
		}(),
	}

	err := s.chatRoomRepo.CreateChatRoom(&chatRoom)
	if err != nil {
		global.Gzap.Error("CreateChatRoom", zap.Error(err))
		return gpt2.ChatRoom{}, errors.New("系统内部错误, 请联系管理员")
	}

	return chatRoom, nil
}
func (s *chatRoomService) CreateChatRoom() (gpt2.ChatRoom, error) {
	chatRoom := gpt2.ChatRoom{}

	if err := s.chatRoomRepo.CreateChatRoom(&chatRoom); err != nil {
		global.Gzap.Error("CreateChatRoom", zap.Error(err))
		return gpt2.ChatRoom{}, errors.New("新建聊天失败, 系统内部错误, 请联系管理员")
	}

	return chatRoom, nil
}

func NewChatRoomService() ChatRoomService {
	return &chatRoomService{chatRoomRepo: repository.NewChatRoomRepository()}
}
