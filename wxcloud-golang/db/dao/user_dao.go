package dao

import (
	"context"
	"fmt"
	"time"
	"wxcloud-golang/db"
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
		return conn.Model(&model.Family{}).
			Select(`
				family.*,
				(SELECT count(1) FROM family_member fm WHERE fm.family_id = family.id) as member_count,
				family_member.role as my_role,
				family_member.sort_order as my_sort_order,
				family_member.created_at as join_time
			`).
			Joins("INNER JOIN family_member ON family_member.family_id = family.id").
			Where("family_member.open_id = ?", OPENID).
			Order("family_member.sort_order ASC, family.id ASC").
			Find(&family).Error
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

func DeleteFamilyWithData(ctx context.Context, familyID uint) error {
	// 开启事务
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("family_id = ?", familyID).Delete(&model.Plant{}).Error; err != nil {
			return err
		}
		if err := tx.Where("family_id = ?", familyID).Delete(&model.Tag{}).Error; err != nil {
			return err
		}
		// 3. (如果有) 删除该家庭下的所有 图片记录/生长记录
		if err := tx.Where("family_id = ?", familyID).Delete(&model.Image{}).Error; err != nil {
			return err
		}
		// 4. 删除 家庭与用户的关联 (如果有中间表 user_families)
		if err := tx.Table("family_member").Where("family_id = ?", familyID).Delete(nil).Error; err != nil {
			return err
		}

		// 5. 最后删除 家庭 (Family) 本身
		if err := tx.Where("id = ?", familyID).Delete(&model.Family{}).Error; err != nil {
			return err
		}

		// 事务提交
		return nil
	})
}

func UpdateFamilySortOrder(ctx context.Context, familyIDs []uint, currentOpenID string) error {
	return execWithSpan(ctx, "UPDATE", "family_member", func(conn *gorm.DB) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			for index, familyID := range familyIDs {
				// SQL 逻辑:
				// UPDATE family_member
				// SET sort_order = index
				// WHERE family_id = familyID AND open_id = currentOpenID

				if err := tx.Model(&model.FamilyMember{}).
					Where("family_id = ? AND open_id = ?", familyID, currentOpenID).
					Update("sort_order", index).Error; err != nil {
					return err
				}
			}
			return nil
		})
	})
}

func SwitchCurrentFamily(ctx context.Context, openID string, familyID uint) error {
	return execWithSpan(ctx, "UPDATE", "user", func(conn *gorm.DB) error {
		// 1. 安全检查：确认用户是该家庭成员
		var count int64
		err := conn.Model(&model.FamilyMember{}).
			Where("family_id = ? AND open_id = ?", familyID, openID).
			Count(&count).Error

		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("你不是该家庭成员，无法切换")
		}

		// 2. 更新用户表的 CurrentFamilyID
		return conn.Model(&model.User{}).
			Where("open_id = ?", openID).
			Update("current_family_id", familyID).Error
	})
}
