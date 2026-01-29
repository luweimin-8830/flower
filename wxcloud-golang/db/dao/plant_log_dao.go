package dao

import (
	"context"
	"time"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// BatchCreatePlantLogs 批量插入日志（包含图片关联逻辑）
func BatchCreatePlantLogs(ctx context.Context, logs []*model.PlantLog, imageIDs []uint) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查图片
		var images []model.Image
		if len(imageIDs) > 0 {
			if err := tx.Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
				return err
			}
		}

		// 2. 内存中关联图片
		if len(images) > 0 {
			for _, log := range logs {
				log.Images = images
			}
		}

		// 3. 批量入库
		if err := tx.CreateInBatches(logs, 100).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPlantLogsByPlantID 查询指定植物的日志
func GetPlantLogsByPlantID(ctx context.Context, plantID uint) ([]*model.PlantLog, error) {
	var logs []*model.PlantLog
	err := db.DB.WithContext(ctx).
		Preload("Images").
		Where("plant_id = ?", plantID).
		Order("log_time desc").
		Find(&logs).Error
	return logs, err
}

// DeletePlantLogByID 删除日志
func DeletePlantLogByID(ctx context.Context, logID uint) error {
	return db.DB.WithContext(ctx).Select("Images").Delete(&model.PlantLog{}, logID).Error
}

// CheckTodayActionExists 检查今日是否已存在某种操作
func CheckTodayActionExists(ctx context.Context, plantID uint, actionType string) (bool, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := db.DB.WithContext(ctx).Model(&model.PlantLog{}).
		Where("plant_id = ? AND action_type = ? AND DATE(log_time) = ?", plantID, actionType, today).
		Count(&count).Error
	return count > 0, err
}
