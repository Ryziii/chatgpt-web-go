package repository

import (
	"chatgpt-web-go/src/global"
	"chatgpt-web-go/src/model/api/gpt"
	"gorm.io/gorm"
)

type ChatMessageRepository interface {
	GetOne(result *gpt.ChatMessageDO, source gpt.ChatMessageDO) error
	CreateChatMessage(chatMessage *gpt.ChatMessageDO) error
	UpdateChatMessage(chatMessage *gpt.ChatMessageDO) error
	DeleteChatMessage(chatMessage *gpt.ChatMessageDO) error
	GetChatMessageByID(result *gpt.ChatMessageDO, id uint) error
}

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository() ChatMessageRepository {
	return &chatMessageRepository{db: global.Gdb.Model(&gpt.ChatMessageDO{})}
}

func (r *chatMessageRepository) GetOne(result *gpt.ChatMessageDO, source gpt.ChatMessageDO) error {
	return r.db.Where(source).First(result).Error
}

func (r *chatMessageRepository) CreateChatMessage(chatMessage *gpt.ChatMessageDO) error {
	return global.Gdb.Create(chatMessage).Error
}

func (r *chatMessageRepository) UpdateChatMessage(chatMessage *gpt.ChatMessageDO) error {
	return global.Gdb.Model(chatMessage).Updates(chatMessage).Error
}

func (r *chatMessageRepository) DeleteChatMessage(chatMessage *gpt.ChatMessageDO) error {
	return global.Gdb.Delete(chatMessage).Error
}

func (r *chatMessageRepository) GetChatMessageByID(chatMessage *gpt.ChatMessageDO, id uint) (err error) {
	if err := global.Gdb.Where("id = ?", id).First(chatMessage).Error; err != nil {
		return err
	}
	return nil
}
