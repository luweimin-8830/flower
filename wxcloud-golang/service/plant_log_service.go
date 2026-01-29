package service

import (
	"context"
	"errors"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

// BatchAddPlantLogs 业务逻辑：添加日志
func BatchAddPlantLogs(ctx context.Context, logs []*model.PlantLog, imageIDs []uint) error {
	// 1. 判断是否今日已有相同操作 (仅针对批量添加中的第一个植物做演示，或循环检查)
	// 如果是快捷护理，通常是针对单棵植物或一批植物的今日去重
	for _, log := range logs {
		exists, err := dao.CheckTodayActionExists(ctx, log.PlantID, log.ActionType)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("今日已记录过该操作")
		}
	}

	// 调用 DAO 层进行存储
	return dao.BatchCreatePlantLogs(ctx, logs, imageIDs)
}

// GetPlantLogs 业务逻辑：获取日志
func GetPlantLogs(ctx context.Context, plantID uint) ([]*model.PlantLog, error) {
	return dao.GetPlantLogsByPlantID(ctx, plantID)
}

// DeletePlantLog 业务逻辑：删除日志
func DeletePlantLog(ctx context.Context, logID uint) error {
	// 比如：检查是否有权限删除（如果需要）
	return dao.DeletePlantLogByID(ctx, logID)
}
