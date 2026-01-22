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

func UpdateSortOrder(ctx context.Context, tagIDs []uint) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for index, id := range tagIDs {
			if err := tx.Model(&model.Tag{}).
				Where("id = ?", id).
				Update("sort_order", index).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func GetTagByFamilyID(ctx context.Context, familyID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := db.DB.WithContext(ctx).
		Table("tag").
		Select("tag.*, count(plant.id) as plant_count").
		Joins("LEFT JOIN plant_tags ON plant_tags.tag_id = tag.id").
		Joins("LEFT JOIN plant ON plant.id = plant_tags.plant_id AND plant.deleted_at IS NULL").
		Where("tag.family_id = ? AND tag.deleted_at IS NULL", familyID).
		Order("tag.sort_order ASC, tag.id ASC").
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
