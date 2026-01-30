package dao

import (
	"context"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

func CreateCareAction(ctx context.Context, care *model.CareAction) error {
	return execWithSpan(ctx, "INSERT", "care_action", func(conn *gorm.DB) error {
		return conn.Create(care).Error
	})
}

func DeleteCareAction(ctx context.Context, id uint) error {
	return execWithSpan(ctx, "DELETE", "care_action", func(conn *gorm.DB) error {
		return conn.Delete(&model.CareAction{}, id).Error
	})
}

func UpdateCareAction(ctx context.Context, care *model.CareAction) error {
	return execWithSpan(ctx, "UPDATE", "care_action", func(conn *gorm.DB) error {
		// 使用 Updates 而不是 Save，避免触发零值日期 (0001-01-01) 的更新导致 MySQL 报错
		return conn.Model(care).Updates(care).Error
	})
}

func UpdateCareSortOrder(ctx context.Context, careIDs []uint) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for index, id := range careIDs {
			if err := tx.Model(&model.CareAction{}).
				Where("id = ?", id).
				Update("sort_order", index).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func GetCareByFamilyID(ctx context.Context, familyID uint) ([]model.CareAction, error) {
	var cares []model.CareAction
	err := execWithSpan(ctx, "SELECT", "care_action", func(conn *gorm.DB) error {
		return conn.Model(&model.CareAction{}).
			Where("family_id = ?", familyID).
			Order("sort_order ASC, id ASC").
			Find(&cares).Error
	})

	return cares, err
}

func GetCareByID(ctx context.Context, id uint) (*model.CareAction, error) {
	var care model.CareAction
	err := execWithSpan(ctx, "SELECT", "care_action", func(conn *gorm.DB) error {
		return conn.First(&care, id).Error
	})
	return &care, err
}

func CreateDefaultCareActions(ctx context.Context, familyID uint) error {
	defaults := []model.CareAction{
		{Name: "浇水", Type: "water", Icon: "plant-jiaoshui1", Color: "#D6EAF8", FamilyID: familyID, SortOrder: 1},
		{Name: "施肥", Type: "fertilize", Icon: "plant-shifei1", Color: "#DCECC9", FamilyID: familyID, SortOrder: 2},
		{Name: "修剪", Type: "prune", Icon: "plant-a-xiujian13", Color: "#F2D7D5", FamilyID: familyID, SortOrder: 3},
		{Name: "换土", Type: "repot", Icon: "plant-a-Frame9", Color: "#E8E0D5", FamilyID: familyID, SortOrder: 4},
	}
	return execWithSpan(ctx, "INSERT", "care_action", func(conn *gorm.DB) error {
		return conn.Create(&defaults).Error
	})
}
