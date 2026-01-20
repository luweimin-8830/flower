package service

import (
	"context"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

// BatchAddPlantLogs 业务逻辑：添加日志
func BatchAddPlantLogs(ctx context.Context, logs []*model.PlantLog, imageIDs []uint) error {
	// 这里可以添加额外的业务逻辑，比如：
	// 1. 检查植物是否存在
	// 2. 如果是浇水操作，更新植物表的 "last_water_date"
	
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
