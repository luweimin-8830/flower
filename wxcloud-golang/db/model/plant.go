package model

import (
	"time"

	"gorm.io/gorm"
)

type Plant struct {
	gorm.Model
	Name     string    `json:"name" gorm:"type:varchar(100);not null"` //名称
	Cover    string    `json:"cover" gorm:"type:varchar(255)"`         //主图
	Desc     string    `json:"desc" gorm:"type:varchar(255)"`          //备注
	Labels   string    `json:"labels" gorm:"type:varchar(255)"`        //分类
	OpenId   string    `json:"openId"`                                 //所属人
	Tags     string    `json:"tags" gorm:"type:varchar(255)"`          //标签
	Birthday time.Time `json:"time"`                                   //到家时间
}

type PlantListReq struct {
	Page     int    `form:"page" json:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" json:"page_size"`
	Labels   string `form:"labels" json:"labels"`
}
