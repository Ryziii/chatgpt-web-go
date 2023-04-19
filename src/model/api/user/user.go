package user

import (
	"chatgpt-web-go/src/model/api/user/setting"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string           `json:"username"`
	Password string           `json:"password"`
	General  setting.General  `json:"general" gorm:"type:text"`
	Usage    setting.Usage    `json:"usage" gorm:"type:text"`
	Advanced setting.Advanced `json:"advanced" gorm:"type:text"`
}
