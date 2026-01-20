package dao

import (
	"context"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// 创建植物
func CreatePlant(ctx context.Context, plant *model.Plant) error {
	return execWithSpan(ctx, "INSERT", "plant", func(conn *gorm.DB) error {
		return conn.Create(plant).Error
	})
}

func GetPlantByID(ctx context.Context, id uint) (*model.Plant, error) {
	var plant model.Plant
	err := execWithSpan(ctx, "SELECT", "plant", func(conn *gorm.DB) error {
		return conn.Preload("Cover").Preload("Tags").First(&plant, id).Error
	})
	return &plant, err
}

func GetPlantByFamilyID(ctx context.Context, familyID uint) ([]model.Plant, error) {
	var plants []model.Plant
	err := execWithSpan(ctx, "SELECT", "plant", func(conn *gorm.DB) error {
		return conn.Where("family_id = ?", familyID).Preload("Cover").Preload("Tags").Find(&plants).Error
	})
	return plants, err
}

func DeletePlant(ctx context.Context, id uint) error {
	return execWithSpan(ctx, "DELETE", "plant", func(conn *gorm.DB) error {
		return conn.Delete(&model.Plant{}, id).Error
	})
}

func UpdatePlant(ctx context.Context, plantID uint, updateDate map[string]any, newTagIDs []uint) error {
	return execWithSpan(ctx, "UPDATE", "plant", func(conn *gorm.DB) error {
		tx := conn.Begin()
		if err := tx.Model(&model.Plant{}).Where("id = ?", plantID).Updates(updateDate).Error; err != nil {
			tx.Rollback()
			return err
		}

		if newTagIDs != nil {
			var plant model.Plant
			plant.ID = plantID

			var newTags []model.Tag
			for _, tagId := range newTagIDs {
				newTags = append(newTags, model.Tag{Model: gorm.Model{ID: tagId}})
			}
			if err := tx.Model(&plant).Association("Tags").Replace(newTags); err != nil {
				tx.Rollback()
				return err
			}
		}
		return tx.Commit().Error
	})
}
