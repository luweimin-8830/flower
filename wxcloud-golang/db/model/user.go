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
	CurrentFamilyID *uint  `json:"currentFamilyId" gorm:"default:null"`
	RemindTime      string `json:"remindTime" gorm:"type:varchar(5);default:'08:00'"`
	Longitude       float64 `json:"longitude" gorm:"type:decimal(10,7)"`
	Latitude        float64 `json:"latitude" gorm:"type:decimal(10,7)"`
	City            string  `json:"city" gorm:"type:varchar(50)"`
	Adcode          string  `json:"adcode" gorm:"type:varchar(20)"`
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
