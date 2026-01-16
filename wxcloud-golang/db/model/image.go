package model

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	URL string `json:"url" gorm:"type:varchar(255);not null"`
	// 核心：Hash 设为唯一索引，物理上保证不会有重复图片记录
	Hash     string `json:"hash" gorm:"type:char(64);uniqueIndex;not null"`
	Width    int    `json:"width"`                            // 宽
	Height   int    `json:"height"`                           // 高
	Size     int64  `json:"size"`                             // 文件大小(字节)
	MimeType string `json:"mimeType" gorm:"type:varchar(50)"` // e.g. image/jpeg
}
