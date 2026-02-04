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
		err := conn.Preload("Cover").Preload("Tags").First(&plant, id).Error
		if err == nil && plant.Cover.URL == "" {
			plant.Cover.URL = "/static/default.svg"
			plant.Cover.Width = 100
			plant.Cover.Height = 100
		}
		return err
	})
	return &plant, err
}

func GetPlantByFamilyID(ctx context.Context, familyID uint) ([]model.Plant, error) {
	var plants []model.Plant
	err := execWithSpan(ctx, "SELECT", "plant", func(conn *gorm.DB) error {
		err := conn.Where("family_id = ? AND death_at IS NULL", familyID).Order("created_at desc").Preload("Cover").Preload("Tags").Find(&plants).Error
		if err == nil {
			for i := range plants {
				if plants[i].Cover.URL == "" {
					plants[i].Cover.URL = "/static/default.svg"
					plants[i].Cover.Width = 100
					plants[i].Cover.Height = 100
				}
			}
		}
		return err
	})
	return plants, err
}

func GetPlantsPaged(ctx context.Context, familyID uint, isDead bool, page, pageSize int, tagID uint, keyword string) ([]model.Plant, int64, error) {
	var plants []model.Plant
	var total int64
	err := execWithSpan(ctx, "SELECT", "plant", func(conn *gorm.DB) error {
		query := conn.Model(&model.Plant{}).Where("family_id = ?", familyID)
		if isDead {
			query = query.Where("death_at IS NOT NULL")
		} else {
			query = query.Where("death_at IS NULL")
		}

		if tagID > 0 {
			query = query.Joins("JOIN plant_tags ON plant_tags.plant_id = plant.id").Where("plant_tags.tag_id = ?", tagID)
		}

		if keyword != "" {
			query = query.Where("name LIKE ?", "%"+keyword+"%")
		}

		query.Count(&total)

		if page > 0 && pageSize > 0 {
			offset := (page - 1) * pageSize
			query = query.Limit(pageSize).Offset(offset)
		}

		err := query.Order("created_at desc").Preload("Cover").Preload("Tags").Find(&plants).Error
		if err == nil {
			for i := range plants {
				if plants[i].Cover.URL == "" {
					plants[i].Cover.URL = "/static/default.svg"
					plants[i].Cover.Width = 100
					plants[i].Cover.Height = 100
				}
			}
		}
		return err
	})
	return plants, total, err
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
				newTags = append(newTags, model.Tag{ID: tagId})
			}
			if err := tx.Model(&plant).Association("Tags").Replace(newTags); err != nil {
				tx.Rollback()
				return err
			}
		}
		return tx.Commit().Error
	})
}
