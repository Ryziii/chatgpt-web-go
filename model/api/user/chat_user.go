package user

import (
	setting2 "chatgpt-web-go/model/api/user/setting"
	"chatgpt-web-go/model/common"
)

type User struct {
	common.Model
	Username string            `json:"username"`
	Password string            `json:"password"`
	General  setting2.General  `json:"general" gorm:"type:text"`
	Usage    setting2.Usage    `json:"usage" gorm:"type:text"`
	Advanced setting2.Advanced `json:"advanced" gorm:"type:text"`
}

func (User) TableName() string {
	return "chat_user"
}
