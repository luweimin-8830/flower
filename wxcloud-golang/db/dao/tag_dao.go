package dao

import (
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

func CreateTag(tag *model.Tag) error {
	return db.DB.Create(tag).Error
}

func DeleteTag(tagID uint) error {
	return db.DB.Delete(&model.Tag{}, tagID).Error
}

func UpdateTag(tagID uint, name string) error {
	return db.DB.Model(&model.Tag{}).Where("id = ?", tagID).Update("name", name).Error
}

func GetTagByFamilyID(familyID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := db.DB.Where("family_id = ?", familyID).Find(&tags).Error
	return tags, err
}

func GetTagByID(tagID uint) (*model.Tag, error) {
	var tag model.Tag
	err := db.DB.First(&tag, tagID).Error
	return &tag, err
}
