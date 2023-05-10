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
	GetChatMessageById(result *gpt.ChatMessage, id uint64) error
}

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository() ChatMessageRepository {
	return &chatMessageRepository{db: global.Gdb}
}

func (r *chatMessageRepository) GetOne(result *gpt.ChatMessage, source gpt.ChatMessage) error {
	return r.db.Model(gpt.ChatMessage{}).Where(source).First(result).Error
}

func (r *chatMessageRepository) CreateChatMessage(chatMessage *gpt.ChatMessage) error {
	return r.db.Model(gpt.ChatMessage{}).Create(chatMessage).Error
}

func (r *chatMessageRepository) UpdateChatMessage(chatMessage *gpt.ChatMessage) error {
	return r.db.Model(chatMessage).Updates(chatMessage).Error
}

func (r *chatMessageRepository) DeleteChatMessage(chatMessage *gpt.ChatMessage) error {
	return r.db.Model(gpt.ChatMessage{}).Delete(chatMessage).Error
}

func (r *chatMessageRepository) GetChatMessageById(chatMessage *gpt.ChatMessage, id uint64) (err error) {
	if err := r.db.Model(gpt.ChatMessage{}).Where("id = ?", id).First(chatMessage).Error; err != nil {
		return err
	}
	return nil
}
