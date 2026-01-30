package model

import (
	"time"

	"gorm.io/gorm"
)

// 标签表
type Tag struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `json:"name" gorm:"type:varchar(50);not null"`
	FamilyID   uint           `json:"familyId" gorm:"index;not null"`
	Color      string         `json:"color"`
	SortOrder  int            `json:"sortOrder" gorm:"default:0"`
	PlantCount int64          `json:"plantCount" gorm:"->"`
}
