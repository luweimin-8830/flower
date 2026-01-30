package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户表
type User struct {
	ID              uint   `json:"id"`
	OPENID          string `json:"openId" gorm:"uniqueIndex;type:varchar(64)"`
	Phone           uint16 `json:"phone"`
	Name            string `json:"name"`
	CreatedAT       time.Time
	LastDateAT      time.Time
	CurrentFamilyID *uint `json:"currentFamilyId" gorm:"default:null"`
}

type Family struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name" gorm:"type:varchar(50)"`
	OwnerOpenId string         `json:"ownerOpenId" gorm:"index"`
	MemberCount int64          `json:"memberCount" gorm:"->"`
	MyRole      string         `json:"myRole" gorm:"->"`
	MySortOrder int            `json:"mySortOrder" gorm:"->"`
	JoinTime    time.Time      `json:"joinTime" gorm:"->"`
}

type FamilyMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyID  uint           `json:"familyId" gorm:"index"`
	OpenID    string         `json:"openId" gorm:"index"`
	Role      string         `json:"role" gorm:"type:varchar(20)"`
	SortOrder int            `json:"sortOrder" gorm:"default:0"`
}
