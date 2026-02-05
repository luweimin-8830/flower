package model

import (
	"time"
)

// WeatherCache 天气缓存表
type WeatherCache struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	City        string    `json:"city" gorm:"index;type:varchar(50)"`
	Adcode      string    `json:"adcode" gorm:"index;type:varchar(20)"`
	Date        string    `json:"date" gorm:"index;type:varchar(10)"` // 格式：2024-03-20
	Temperature string    `json:"temperature" gorm:"type:varchar(10)"`
	Weather     string    `json:"weather" gorm:"type:varchar(20)"`
	Icon        string    `json:"icon" gorm:"type:varchar(50)"`
	Casts       string    `json:"casts" gorm:"type:text"` // 存储全量预报 JSON 字符串
}
