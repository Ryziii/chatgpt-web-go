package repository

import (
	"chatgpt-web-go/global"
	"chatgpt-web-go/model/api/gpt"
	"gorm.io/gorm"
)

type ChatRoomRepository interface {
	GetOne(result *gpt.ChatRoom, source gpt.ChatRoom) error
	CreateChatRoom(chatRoom *gpt.ChatRoom) error
	UpdateChatRoom(chatRoom *gpt.ChatRoom) error
}

func NewChatRoomRepository() ChatRoomRepository {
	return &chatRoomRepository{db: global.Gdb.Model(&gpt.ChatRoom{})}
}

type chatRoomRepository struct {
	db *gorm.DB
}

func (r *chatRoomRepository) CreateChatRoom(chatRoom *gpt.ChatRoom) error {
	return r.db.Create(chatRoom).Error
}

func (r *chatRoomRepository) UpdateChatRoom(chatRoom *gpt.ChatRoom) error {
	return r.db.Model(chatRoom).Updates(chatRoom).Error
}

func (r *chatRoomRepository) GetOne(result *gpt.ChatRoom, source gpt.ChatRoom) error {
	return r.db.Where(source).First(result).Error
}
