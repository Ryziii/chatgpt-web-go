package repository

import (
	"chatgpt-web-go/src/global"
	"chatgpt-web-go/src/model/api/gpt"
	"gorm.io/gorm"
)

type ChatRoomRepository interface {
	GetOne(result *gpt.ChatRoomDO, source gpt.ChatRoomDO) error
	CreateChatRoom(chatRoom *gpt.ChatRoomDO) error
	UpdateChatRoom(chatRoom *gpt.ChatRoomDO) error
}

func NewChatRoomRepository() ChatRoomRepository {
	return &chatRoomRepository{db: global.Gdb.Model(&gpt.ChatRoomDO{})}
}

type chatRoomRepository struct {
	db *gorm.DB
}

func (r *chatRoomRepository) CreateChatRoom(chatRoom *gpt.ChatRoomDO) error {
	return r.db.Create(chatRoom).Error
}

func (r *chatRoomRepository) UpdateChatRoom(chatRoom *gpt.ChatRoomDO) error {
	return r.db.Model(chatRoom).Updates(chatRoom).Error
}

func (r *chatRoomRepository) GetOne(result *gpt.ChatRoomDO, source gpt.ChatRoomDO) error {
	return r.db.Where(source).First(result).Error
}
