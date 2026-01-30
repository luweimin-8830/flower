package service

import (
	"context"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

func GetFamilyCareActions(ctx context.Context, familyID uint) ([]model.CareAction, error) {
	cares, err := dao.GetCareByFamilyID(ctx, familyID)
	if err != nil {
		return nil, err
	}

	// 如果没有任何养护项，则初始化默认项 (兼容老旧家庭数据)
	if len(cares) == 0 && familyID > 0 {
		err = dao.CreateDefaultCareActions(ctx, familyID)
		if err != nil {
			return nil, err
		}
		// 重新获取
		return dao.GetCareByFamilyID(ctx, familyID)
	}

	return cares, nil
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
