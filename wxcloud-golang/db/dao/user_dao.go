package dao

import (
	"context"
	"time"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// 创建用户
func CreateUser(ctx context.Context, user *model.User) error {
	return execWithSpan(ctx, "INSERT", "user", func(conn *gorm.DB) error {
		return conn.Create(user).Error
	})
}

// 获取用户列表
func GetUserByOpenID(ctx context.Context, OPENID string) (*model.User, error) {
	var user model.User
	err := execWithSpan(ctx, "SELECT", "user", func(conn *gorm.DB) error {
		return conn.Where("open_id = ?", OPENID).First(&user).Error
	})
	return &user, err
}

// 更新用户列表
func UpdateUser(ctx context.Context, user *model.User) error {
	return execWithSpan(ctx, "UPDATE", "user", func(conn *gorm.DB) error {
		return conn.Save(user).Error
	})
}

// 更新最后登录时间
func UpdateUserLastLogin(ctx context.Context, userID uint) error {
	return execWithSpan(ctx, "UPDATE", "user", func(conn *gorm.DB) error {
		return conn.Model(&model.User{}).Where("id = ?", userID).Update("last_date_at", time.Now()).Error
	})
}

// DeleteUser 删除用户
func DeleteUser(ctx context.Context, id int) error {
	return execWithSpan(ctx, "DELETE", "user", func(conn *gorm.DB) error {
		return conn.Delete(&model.User{}, id).Error
	})
}

// 创建家庭
func CreateFamily(ctx context.Context, family *model.Family) error {
	return execWithSpan(ctx, "INSERT", "family", func(conn *gorm.DB) error {
		return conn.Create(family).Error
	})
}

// 创建家庭成员
func CreateFamilyMember(ctx context.Context, familyMember *model.FamilyMember) error {
	return execWithSpan(ctx, "INSERT", "family_member", func(conn *gorm.DB) error {
		return conn.Create(familyMember).Error
	})
}

// 查询家庭列表
func GetFamilyList(ctx context.Context, OPENID string) ([]model.Family, error) {
	var family []model.Family
	err := execWithSpan(ctx, "SELECT", "family", func(conn *gorm.DB) error {
		return conn.Model(&model.Family{}).Joins("INNER JOIN family_member ON family_member.family_id = family.id").
			Where("family_member.open_id = ?", OPENID).Find(&family).Error
	})
	return family, err
}

// 查询是否存在家庭
func GetFamilyByOpenId(ctx context.Context, OPENID string) (*model.Family, error) {
	var family model.Family
	err := execWithSpan(ctx, "SELECT", "family", func(conn *gorm.DB) error {
		return conn.Model(&model.Family{}).Joins("INNER JOIN family_member ON family_member.family_id = family.id").
			Where("family_member.open_id = ?", OPENID).First(&family).Error
	})
	return &family, err
}
