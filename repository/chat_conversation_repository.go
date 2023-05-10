package repository

import (
	"chatgpt-web-go/global"
	"chatgpt-web-go/model/api/gpt"
	"errors"
	"go.uber.org/zap"
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
	return &chatConversationRepository{db: global.Gdb}
}

func (r *chatConversationRepository) GetOne(result *gpt.ChatConversation, source gpt.ChatConversation) error {
	var res = gpt.ChatConversation{}
	if err := r.db.Model(gpt.ChatConversation{}).
		Preload("Question").
		Preload("Answer").
		Preload("ChatRoom").
		Where(source).First(&res).Error; err != nil {
		return err
	}
	*result = res
	return nil
}

func (r *chatConversationRepository) CreateChatConversation(chatConversation *gpt.ChatConversation) error {
	return r.db.Model(gpt.ChatConversation{}).Create(chatConversation).Error
}

func (r *chatConversationRepository) UpdateChatConversation(chatConversation *gpt.ChatConversation) error {
	return r.db.Model(chatConversation).Updates(chatConversation).Error
}

func (r *chatConversationRepository) DeleteChatConversation(chatConversation *gpt.ChatConversation) error {
	return r.db.Model(gpt.ChatConversation{}).Delete(chatConversation).Error
}

func (r *chatConversationRepository) GetChatConversationById(chatConversation *gpt.ChatConversation, id uint64) error {
	result := &gpt.ChatConversation{}
	query := r.db.Model(gpt.ChatConversation{}).Preload("Question").Preload("Answer").Preload("ChatRoom")
	query = query.Where("id = ?", id)

	if err := query.First(result).Error; err != nil {
		global.Gzap.Error("GetChatConversationById error", zap.Error(err))
		return errors.New("系统内部错误, 请联系管理员")
	}
	*chatConversation = *result
	return nil
}

func (r *chatConversationRepository) GetChatConversationByQuery(chatConversation *gpt.ChatConversation, query map[string]interface{}) error {
	var result gpt.ChatConversation
	if err := r.db.Model(gpt.ChatConversation{}).
		Preload("Question").
		Preload("Answer").
		Preload("ChatRoom").
		Where(query).First(&result).Error; err != nil {
		return err
	}
	*chatConversation = result
	return nil
}
