package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	LogTypeWater     = "water"     // 浇水
	LogTypeFertilize = "fertilize" // 施肥
	LogTypePrune     = "prune"     // 修剪
	LogTypeRepot     = "repot"     // 换盆
	LogTypeRecord    = "record"    // 日常记录
	LogTypePest      = "pest"      // 虫害
)

type PlantLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	PlantID    uint           `json:"plantId" gorm:"index;not null"`               // 关联植物ID
	OpenId     string         `json:"openId" gorm:"index"`                         // 操作人
	ActionType string         `json:"actionType" gorm:"type:varchar(50);not null"` // 类型
	Content    string         `json:"content" gorm:"type:text"`                    // 文字内容
	LogTime    time.Time      `json:"logTime" gorm:"index"`                        // 实际发生的日期

	// 冗余展示字段（不存数据库）
	ActionName string `json:"name" gorm:"-"`
	ActionIcon string `json:"icon" gorm:"-"`
	ActionColor string `json:"color" gorm:"-"`

	// 多对多关联：一条日志可以包含多张图片
	// gorm 会自动创建 plant_log_images 中间表
	Images []Image `json:"images" gorm:"many2many:plant_log_images;"`
}

// GetActionName 获取动作类型的中文名称
func GetActionName(actionType string) string {
	switch actionType {
	case LogTypeWater:
		return "浇水"
	case LogTypeFertilize:
		return "施肥"
	case LogTypePrune:
		return "修剪"
	case LogTypeRepot:
		return "换盆"
	case LogTypeRecord:
		return "日常记录"
	case LogTypePest:
		return "虫害"
	default:
		return actionType
	}
}
