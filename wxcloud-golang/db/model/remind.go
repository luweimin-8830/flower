package model

import (
	"time"
)

type RemindTask struct {
	ID         uint      `gorm:"primaryKey"`
	OpenID     string    `gorm:"type:varchar(128);index"`
	FamilyID   uint      `gorm:"index"`
	PlantID    uint      `gorm:"index"`
	RemindTime time.Time `gorm:"index"`
	Content    string    `gorm:"type:varchar(255)"`
	ActionType string    `gorm:"type:varchar(50)"`
	Status     int       `gorm:"default:0"` // 0: 待发送, 1: 已发送, 2: 发送失败
	CreatedAt  time.Time
}
