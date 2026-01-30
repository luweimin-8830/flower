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
	logs, err := dao.GetPlantLogsByPlantID(ctx, plantID)
	if err != nil {
		return nil, err
	}
	// 获取植物详情以拿到 familyID
	plant, _ := dao.GetPlantByID(ctx, plantID)
	if plant != nil {
		fillActionInfo(ctx, plant.FamilyID, logs)
	}
	return logs, nil
}

// GetFamilyLogs 获取家庭下所有植物的日志
func GetFamilyLogs(ctx context.Context, familyID uint) ([]*model.PlantLog, error) {
	logs, err := dao.GetPlantLogsByFamilyID(ctx, familyID)
	if err != nil {
		return nil, err
	}
	fillActionInfo(ctx, familyID, logs)
	return logs, nil
}

// fillActionInfo 填充护理操作的详细信息
func fillActionInfo(ctx context.Context, familyID uint, logs []*model.PlantLog) {
	if len(logs) == 0 {
		return
	}

	// 获取该家庭的所有养护项配置
	cares, _ := dao.GetCareByFamilyID(ctx, familyID)
	careMap := make(map[string]model.CareAction)
	for _, c := range cares {
		careMap[c.Type] = c
	}

	for _, log := range logs {
		if action, ok := careMap[log.ActionType]; ok {
			log.ActionName = action.Name
			log.ActionIcon = action.Icon
			log.ActionColor = action.Color
		} else {
			// 兜底显示
			log.ActionName = log.ActionType
		}
	}
}

// DeletePlantLog 业务逻辑：删除日志
func DeletePlantLog(ctx context.Context, logID uint) error {
	// 比如：检查是否有权限删除（如果需要）
	return dao.DeletePlantLogByID(ctx, logID)
}

// UpdatePlantLog 业务逻辑：更新日志
func UpdatePlantLog(ctx context.Context, log *model.PlantLog, imageIDs []uint) error {
	return dao.UpdatePlantLog(ctx, log, imageIDs)
}

// GetPlantLog 业务逻辑：获取单条日志
func GetPlantLog(ctx context.Context, logID uint) (*model.PlantLog, error) {
	log, err := dao.GetPlantLogByID(ctx, logID)
	if err != nil {
		return nil, err
	}
	// 获取植物详情以拿到 familyID
	plant, _ := dao.GetPlantByID(ctx, log.PlantID)
	if plant != nil {
		fillActionInfo(ctx, plant.FamilyID, []*model.PlantLog{log})
	}
	return log, nil
}
