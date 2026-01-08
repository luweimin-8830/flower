package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户表
type User struct {
	ID         uint   `json:"id"`
	OPENID     string `json:"openId" gorm:"uniqueIndex;type:varchar(64)"`
	Phone      uint16 `json:"phone"`
	Name       string `json:"name"`
	CreatedAT  time.Time
	LastDateAT time.Time
}

type Family struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(50)"`
	OwnerOpenId string `json:"ownerOpenId" gorm:"index"`
}

type FamilyMember struct {
	gorm.Model
	FamilyID uint   `json:"familyId" gorm:"index"`
	OpenID   string `json:"openId" gorm:"index"`
	Role     string `json:"role" gorm:"type:varchar(20)"`
	// owner: 最高权限，可以解散家庭
	// admin: 可以添加/删除植物，邀请成员
	// member: 只能查看/浇水
}
