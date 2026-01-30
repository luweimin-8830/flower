package service

import (
	"context"
	"fmt"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

func AddPlant(ctx context.Context, plant *model.Plant, tagIDs []uint) error {
	if len(tagIDs) > 0 {
		var tags []model.Tag
		for _, id := range tagIDs {
			tags = append(tags, model.Tag{ID: id})
		}
		plant.Tags = tags
		fmt.Println("传入标签是：", tags)
	}
	return dao.CreatePlant(ctx, plant)
}

func UpdatePlant(ctx context.Context, plantID uint, updateDate map[string]any, tagIDs []uint) error {
	return dao.UpdatePlant(ctx, plantID, updateDate, tagIDs)
}

func GetPlants(ctx context.Context, familyID uint) ([]model.Plant, error) {
	return dao.GetPlantByFamilyID(ctx, familyID)
}

func GetPlant(ctx context.Context, id uint) (*model.Plant, error) {
	return dao.GetPlantByID(ctx, id)
}

func DeletePlant(ctx context.Context, plantID uint) error {
	return dao.DeletePlant(ctx, plantID)
}
