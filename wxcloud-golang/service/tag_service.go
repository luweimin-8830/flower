package service

import (
	"context"
	"errors"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

func AddTag(ctx context.Context, name string, familyID uint) (*model.Tag, error) {
	existingTags, _ := dao.GetTagByFamilyID(ctx, familyID)
	for _, t := range existingTags {
		if t.Name == name {
			return nil, errors.New("家庭已存在此标签名")
		}
	}
	tag := &model.Tag{
		Name:     name,
		FamilyID: familyID,
	}
	err := dao.CreateTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func DeleteTag(ctx context.Context, tagID uint) error {
	return dao.DeleteTag(ctx, tagID)
}

func UpdateTag(ctx context.Context, tagID uint, name string) error {
	return dao.UpdateTag(ctx, tagID, name)
}

func SortTag(ctx context.Context, tagIDs []uint) error {
	return dao.UpdateSortOrder(ctx, tagIDs)
}

func GetFamilyTag(ctx context.Context, familyID uint) ([]model.Tag, error) {
	return dao.GetTagByFamilyID(ctx, familyID)
}
