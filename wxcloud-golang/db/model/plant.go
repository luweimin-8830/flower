package model

import (
	"time"

	"gorm.io/gorm"
)

type Plant struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"` //名称
	CoverID   uint           `json:"coverId"`                                //主图
	Cover     Image          `json:"cover" gorm:"foreignKey:CoverID"`
	Desc      string         `json:"desc" gorm:"type:varchar(255)"` //备注
	OpenId    string         `json:"openId"`                        //所属人
	Birthday  time.Time      `json:"birthday"`                      //到家时间
	DeathAt   *time.Time     `json:"deathAt" gorm:"index"`          //死亡时间
	FamilyID  uint           `json:"familyId" gorm:"index"`
	Tags      []Tag          `json:"tags" gorm:"many2many:plant_tags;"`
	DaysCount int64          `json:"daysCount" gorm:"-"`
}
