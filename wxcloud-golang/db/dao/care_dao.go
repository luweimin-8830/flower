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
		return conn.Save(care).Error
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
	err := db.DB.WithContext(ctx).
		Where("family_id = ?", familyID).
		Order("sort_order ASC, id ASC").
		Find(&cares).Error

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
		{Name: "浇水", Type: "Watering", Icon: "checkbox-filled", Color: "#D6EAF8", FamilyID: familyID, SortOrder: 0},
		{Name: "施肥", Type: "Fertilizing", Icon: "flask", Color: "#DCECC9", FamilyID: familyID, SortOrder: 1},
		{Name: "修剪", Type: "Pruning", Icon: "scissors", Color: "#F2D7D5", FamilyID: familyID, SortOrder: 2},
		{Name: "换土", Type: "SoilChange", Icon: "download", Color: "#E8E0D5", FamilyID: familyID, SortOrder: 3},
	}
	return execWithSpan(ctx, "INSERT", "care_action", func(conn *gorm.DB) error {
		return conn.Create(&defaults).Error
	})
}
