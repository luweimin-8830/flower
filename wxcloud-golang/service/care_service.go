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

	// 针对老用户：如果列表不为空但缺少“成长记录” (Type: "Growth")，补齐它
	hasGrowth := false
	for _, c := range cares {
		if c.Type == "Growth" {
			hasGrowth = true
			break
		}
	}
	if !hasGrowth && familyID > 0 {
		growth := model.CareAction{
			Name:      "成长记录",
			Type:      "Growth",
			Icon:      "camera",
			Color:     "#E8E0D5",
			FamilyID:  familyID,
			SortOrder: len(cares),
		}
		_ = dao.CreateCareAction(ctx, &growth)
		// 重新获取完整列表
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
