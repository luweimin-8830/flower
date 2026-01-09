package dao

import (
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// 创建植物
func CreatePlant(plant *model.Plant) error {
	return db.DB.Create(plant).Error
}

func GetPlantByID(id uint) (*model.Plant, error) {
	var plant model.Plant
	err := db.DB.Preload("Tags").First(&plant, id).Error
	return &plant, err
}

func GetPlantByFamilyID(familyID uint) ([]model.Plant, error) {
	var plants []model.Plant
	err := db.DB.Where("family_id = ?", familyID).Preload("Tags").Find(&plants).Error
	return plants, err
}

func DeletePlant(id uint) error {
	return db.DB.Delete(&model.Plant{}, id).Error
}

func UpdatePlant(plantID uint, updateDate map[string]interface{}, newTagIDs []uint) error {
	tx := db.DB.Begin()
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
}
