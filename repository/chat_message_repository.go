package repository

import (
	"chatgpt-web-go/global"
	"chatgpt-web-go/model/api/gpt"
	"gorm.io/gorm"
)

type ChatMessageRepository interface {
	GetOne(result *gpt.ChatMessage, source gpt.ChatMessage) error
	CreateChatMessage(chatMessage *gpt.ChatMessage) error
	UpdateChatMessage(chatMessage *gpt.ChatMessage) error
	DeleteChatMessage(chatMessage *gpt.ChatMessage) error
	GetChatMessageById(result *gpt.ChatMessage, id uint) error
}

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository() ChatMessageRepository {
	return &chatMessageRepository{db: global.Gdb.Model(&gpt.ChatMessage{})}
}

func (r *chatMessageRepository) GetOne(result *gpt.ChatMessage, source gpt.ChatMessage) error {
	return r.db.Where(source).First(result).Error
}

func (r *chatMessageRepository) CreateChatMessage(chatMessage *gpt.ChatMessage) error {
	return global.Gdb.Create(chatMessage).Error
}

func (r *chatMessageRepository) UpdateChatMessage(chatMessage *gpt.ChatMessage) error {
	return global.Gdb.Model(chatMessage).Updates(chatMessage).Error
}

func (r *chatMessageRepository) DeleteChatMessage(chatMessage *gpt.ChatMessage) error {
	return global.Gdb.Delete(chatMessage).Error
}

func (r *chatMessageRepository) GetChatMessageById(chatMessage *gpt.ChatMessage, id uint) (err error) {
	if err := global.Gdb.Where("id = ?", id).First(chatMessage).Error; err != nil {
		return err
	}
	return nil
}
