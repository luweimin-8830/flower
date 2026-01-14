package dao

import (
	"context"
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
	err := execWithSpan(ctx, "SELECT", "tag", func(conn *gorm.DB) error {
		return conn.Where("family_id = ?", familyID).Find(&tags).Error
	})
	return tags, err
}

func GetTagByID(ctx context.Context, tagID uint) (*model.Tag, error) {
	var tag model.Tag
	err := execWithSpan(ctx, "SELECT", "tag", func(conn *gorm.DB) error {
		return conn.First(&tag, tagID).Error
	})
	return &tag, err
}
