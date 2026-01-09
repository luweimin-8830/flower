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
	OpenId   string    `json:"openId"`                                 //所属人
	Birthday time.Time `json:"time"`                                   //到家时间
	FamilyID uint      `json:"familyId" gorm:"index"`
	Tags     []Tag     `json:"tags" gorm:"many2many:plant_tags;"`
}
