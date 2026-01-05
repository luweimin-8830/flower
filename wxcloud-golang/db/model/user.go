package model

import "time"

// 用户表
type User struct {
	ID         uint   `json:"id"`
	OPENID     string `json:"openId" gorm:"uniqueIndex;type:varchar(64)"`
	Phone      uint16 `json:"phone"`
	Name       string `json:"name"`
	CreatedAT  time.Time
	LastDateAT time.Time
}
