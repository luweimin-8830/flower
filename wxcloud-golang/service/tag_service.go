package service

import (
	"errors"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

func AddTag(name string, familyID uint) (*model.Tag, error) {
	existingTags, _ := dao.GetTagByFamilyID(familyID)
	for _, t := range existingTags {
		if t.Name == name {
			return nil, errors.New("家庭已存在此标签名")
		}
	}
	tag := &model.Tag{
		Name:     name,
		FamilyID: familyID,
	}
	err := dao.CreateTag(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func DeleteTag(tagID uint) error {
	return dao.DeleteTag(tagID)
}

func UpdateTag(tagID uint, name string) error {
	return dao.UpdateTag(tagID, name)
}

func GetFamilyTag(familyID uint) ([]model.Tag, error) {
	return dao.GetTagByFamilyID(familyID)
}
