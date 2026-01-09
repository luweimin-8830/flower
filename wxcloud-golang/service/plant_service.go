package service

import (
	"fmt"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

func AddPlant(plant *model.Plant, tagIDs []uint) error {
	if len(tagIDs) > 0 {
		var tags []model.Tag
		for _, id := range tagIDs {
			tags = append(tags, model.Tag{Model: gorm.Model{ID: id}})
		}
		plant.Tags = tags
		fmt.Println("传入标签是：", tags)
	}
	return dao.CreatePlant(plant)
}

func UpdatePlant(plantID uint, updateDate map[string]interface{}, tagIDs []uint) error {
	return dao.UpdatePlant(plantID, updateDate, tagIDs)
}

func GetPlants(familyID uint) ([]model.Plant, error) {
	return dao.GetPlantByFamilyID(familyID)
}

func GetPlant(id uint) (*model.Plant, error) {
	return dao.GetPlantByID(id)
}

func DeletePlant(plantID uint) error {
	return dao.DeletePlant(plantID)
}
