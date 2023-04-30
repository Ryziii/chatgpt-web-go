package repository

import (
	"chatgpt-web-go/global"
	"chatgpt-web-go/model/api/gpt"
	"gorm.io/gorm"
)

type ChatConversationRepository interface {
	GetOne(result *gpt.ChatConversation, source gpt.ChatConversation) error
	CreateChatConversation(chatConversation *gpt.ChatConversation) error
	UpdateChatConversation(chatConversation *gpt.ChatConversation) error
	DeleteChatConversation(chatConversation *gpt.ChatConversation) error
	GetChatConversationById(chatConversation *gpt.ChatConversation, id uint64) error
	GetChatConversationByQuery(chatConversation *gpt.ChatConversation, query map[string]interface{}) error
}

type chatConversationRepository struct {
	db *gorm.DB
}

func NewChatConversationRepository() ChatConversationRepository {
	return &chatConversationRepository{db: global.Gdb.Model(&gpt.ChatConversation{})}
}

func (r *chatConversationRepository) GetOne(result *gpt.ChatConversation, source gpt.ChatConversation) error {
	return r.db.
		Preload("Question").
		Preload("Answer").
		Preload("ChatRoom").
		Where(source).First(result).Error
}

func (r *chatConversationRepository) CreateChatConversation(chatConversation *gpt.ChatConversation) error {
	return r.db.Create(chatConversation).Error
}

func (r *chatConversationRepository) UpdateChatConversation(chatConversation *gpt.ChatConversation) error {
	return r.db.Updates(chatConversation).Error
}

func (r *chatConversationRepository) DeleteChatConversation(chatConversation *gpt.ChatConversation) error {
	return r.db.Delete(chatConversation).Error
}

func (r *chatConversationRepository) GetChatConversationById(chatConversation *gpt.ChatConversation, id uint64) error {
	if err := r.db.
		Preload("Question").
		Preload("Answer").
		Preload("ChatRoom").
		Where("id = ?", id).
		First(chatConversation).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatConversationRepository) GetChatConversationByQuery(chatConversation *gpt.ChatConversation, query map[string]interface{}) error {
	return r.db.
		Preload("Question").
		Preload("Answer").
		Preload("ChatRoom").
		Where(query).First(chatConversation).Error
}
