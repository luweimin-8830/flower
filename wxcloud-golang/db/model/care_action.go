package model

import (
	"gorm.io/gorm"
)

// CareAction 养护项配置表
type CareAction struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(50);not null"`
	Icon      string `json:"icon" gorm:"type:varchar(50)"`
	Color     string `json:"color" gorm:"type:varchar(20)"`
	Type      string `json:"type" gorm:"type:varchar(50);not null"` // 如 Watering, Fertilizing
	FamilyID  uint   `json:"familyId" gorm:"index;not null"`
	SortOrder int    `json:"sortOrder" gorm:"default:0"`
}
