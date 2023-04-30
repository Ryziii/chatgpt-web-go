package gpt

import (
	"chatgpt-web-go/global"
	model "chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/repository"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
)

type ChatRoomService interface {
	CreateChatRoomByChatMessage(chatMessageDO *model.ChatMessage) (model.ChatRoom, error)
	CreateChatRoom() (model.ChatRoom, error)
}

type chatRoomService struct {
	chatRoomRepo repository.ChatRoomRepository
}

func (s *chatRoomService) CreateChatRoomByChatMessage(chatMessageDO *model.ChatMessage) (model.ChatRoom, error) {
	chatRoom := model.ChatRoom{
		IP:                 "",
		FirstChatMessageId: chatMessageDO.Id,
		FirstMessageId:     uuid.New().String(),
		Title: func() string {
			ru := []rune(chatMessageDO.Content)
			return string(ru[:int(math.Min(float64(len(ru)), 50))])
		}(),
	}

	err := s.chatRoomRepo.CreateChatRoom(&chatRoom)
	if err != nil {
		global.Gzap.Error("CreateChatRoom", zap.Error(err))
		return model.ChatRoom{}, errors.New("系统内部错误, 请联系管理员")
	}

	return chatRoom, nil
}
func (s *chatRoomService) CreateChatRoom() (model.ChatRoom, error) {
	chatRoom := model.ChatRoom{}

	if err := s.chatRoomRepo.CreateChatRoom(&chatRoom); err != nil {
		global.Gzap.Error("CreateChatRoom", zap.Error(err))
		return model.ChatRoom{}, errors.New("新建聊天失败, 系统内部错误, 请联系管理员")
	}

	return chatRoom, nil
}

func NewChatRoomService() ChatRoomService {
	return &chatRoomService{chatRoomRepo: repository.NewChatRoomRepository()}
}
