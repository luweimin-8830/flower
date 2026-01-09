package model

import (
	"gorm.io/gorm"
)

// 标签表
type Tag struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(50);not null"`
	FamilyID uint   `json:"faimilyId" gorm:"index;not null"`
	Color    string `json:"color"`
}
