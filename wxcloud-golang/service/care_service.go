package service

import (
	"context"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

func GetFamilyCareActions(ctx context.Context, familyID uint) ([]model.CareAction, error) {
	return dao.GetCareByFamilyID(ctx, familyID)
}

func AddCareAction(ctx context.Context, care *model.CareAction) error {
	return dao.CreateCareAction(ctx, care)
}

func UpdateCareAction(ctx context.Context, care *model.CareAction) error {
	return dao.UpdateCareAction(ctx, care)
}

func DeleteCareAction(ctx context.Context, id uint) error {
	return dao.DeleteCareAction(ctx, id)
}

func SortCareActions(ctx context.Context, careIDs []uint) error {
	return dao.UpdateCareSortOrder(ctx, careIDs)
}
