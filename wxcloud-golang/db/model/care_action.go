package model

import (
	"time"

	"gorm.io/gorm"
)

// CareAction 养护项配置表
type CareAction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Icon      string         `json:"icon" gorm:"type:varchar(50)"`
	Color     string         `json:"color" gorm:"type:varchar(20)"`
	Type      string         `json:"type" gorm:"type:varchar(50);not null"` // 如 Watering, Fertilizing
	FamilyID  uint           `json:"familyId" gorm:"index;not null"`
	SortOrder int            `json:"sortOrder" gorm:"default:0"`
}
