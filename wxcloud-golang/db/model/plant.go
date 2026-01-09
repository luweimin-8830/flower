package model

import (
	"time"

	"gorm.io/gorm"
)

type Plant struct {
	gorm.Model
	Name            string    `json:"name" gorm:"type:varchar(100);not null"` //名称
	Cover           string    `json:"cover" gorm:"type:varchar(255)"`         //主图
	Desc            string    `json:"desc" gorm:"type:varchar(255)"`          //备注
	OpenId          string    `json:"openId"`                                 //所属人
	Birthday        time.Time `json:"time"`                                   //到家时间
	CurrentUserRole string    `json:"role" gorm:"-"`                          //权限
	FamilyID        uint      `json:"familyId" gorm:"index"`
	Tags            []Tag     `json:"tags" gorm:"many2many:plant_tags;"`
}

// 用户权限表 弃用,待删除
type PlantUserRelation struct {
	gorm.Model
	PlantID uint   `gorm:"index;not null"`
	OpenId  string `gorm:"index;not null"`
	Role    string `gorm:"type:varchar(20);default:'viewer'"` //权限 'admin','editor','viewer' 管理员,可编辑,仅查看
}

type PlantListReq struct {
	Page     int    `form:"page" json:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" json:"page_size"`
	Labels   string `form:"labels" json:"labels"`
}

type PlantAddReq struct {
	Name     string `json:"name" binding:"required"` // 必填
	Cover    string `json:"cover"`
	Desc     string `json:"desc"`
	Labels   string `json:"labels"`
	Tags     string `json:"tags"`
	Birthday string `json:"birthday"`
	FamilyID uint   `json:"familyId"`
}

type PlantDeleteReq struct {
	ID uint `json:"id" binding:"required"`
}

type PlantUpdateReq struct {
	ID       uint     `json:"id" binding:"required"`
	Name     *string  `json:"name"`
	Cover    *string  `json:"cover"`
	Desc     *string  `json:"desc"`
	Labels   []string `json:"labels"`
	Tags     []string `json:"tags"`
	Birthday *string  `json:"birthday"`
}
