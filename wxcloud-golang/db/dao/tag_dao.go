package dao

import (
	"context"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

func CreateTag(ctx context.Context, tag *model.Tag) error {
	return execWithSpan(ctx, "INSERT", "tag", func(conn *gorm.DB) error {
		return conn.Create(tag).Error
	})
}

func DeleteTag(ctx context.Context, tagID uint) error {
	return execWithSpan(ctx, "DELETE", "tag", func(conn *gorm.DB) error {
		return conn.Delete(&model.Tag{}, tagID).Error
	})
}

func UpdateTag(ctx context.Context, tagID uint, name string) error {
	return execWithSpan(ctx, "UPDATE", "tag", func(conn *gorm.DB) error {
		return conn.Model(&model.Tag{}).Where("id = ?", tagID).Update("name", name).Error
	})
}

func GetTagByFamilyID(ctx context.Context, familyID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := db.DB.WithContext(ctx).
		Table("tag").
		Select("tag.*, count(plant_tags.plant_id) as plant_count").
		// 注意：'plant_tags' 是你的多对多中间表表名，请去数据库确认一下是否叫这个
		Joins("LEFT JOIN plant_tags ON plant_tags.tag_id = tag.id").
		Where("tag.family_id = ?", familyID).
		Group("tag.id").
		Scan(&tags).Error
	return tags, err
}

func GetTagByID(ctx context.Context, tagID uint) (*model.Tag, error) {
	var tag model.Tag
	err := execWithSpan(ctx, "SELECT", "tag", func(conn *gorm.DB) error {
		return conn.First(&tag, tagID).Error
	})
	return &tag, err
}
