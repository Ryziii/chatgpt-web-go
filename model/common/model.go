package common

import "time"

type Model struct {
	Id         uint64    `gorm:"primarykey"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}
