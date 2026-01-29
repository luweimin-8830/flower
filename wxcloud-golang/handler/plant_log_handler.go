package handler

import (
	"time"
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

// BatchCreateLogRequest 批量创建请求
type BatchCreateLogRequest struct {
	PlantIDs   []uint `json:"plantIds" binding:"required"` // 支持多选植物 [1, 2, 3]
	ActionType string `json:"actionType" binding:"required"`
	Content    string `json:"content"`
	LogTime    string `json:"logTime"`  // 前端传 "2026-01-20"
	ImageIDs   []uint `json:"imageIds"` // 图片ID列表
}

type GetLogsRequest struct {
	PlantID uint `json:"plantId" binding:"required"`
}

type DeleteLogRequest struct {
	ID uint `json:"id" binding:"required"`
}

type UpdateLogRequest struct {
	ID         uint   `json:"id" binding:"required"`
	ActionType string `json:"actionType"`
	Content    string `json:"content"`
	LogTime    string `json:"logTime"`
	ImageIDs   []uint `json:"imageIds"`
}

type GetSingleLogRequest struct {
	ID uint `json:"id" binding:"required"`
}

// CreatePlantLogHandler 创建日志（支持批量）
func CreatePlantLogHandler(c *gin.Context) {
	var req BatchCreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误: "+err.Error())
		return
	}

	OPENID := c.GetHeader("X-WX-OPENID")

	// 1. 处理时间：如果没传，默认当前时间；传了就解析
	logTime := time.Now()
	if req.LogTime != "" {
		t, err := time.Parse("2006-01-02", req.LogTime)
		if err != nil {
			response.Fail(c, "日期格式错误，应为 YYYY-MM-DD")
			return
		}
		logTime = t
	}

	// 2. 构造日志对象切片
	var logs []*model.PlantLog
	for _, pid := range req.PlantIDs {
		logs = append(logs, &model.PlantLog{
			PlantID:    pid,
			ActionType: req.ActionType,
			Content:    req.Content,
			LogTime:    logTime,
			OpenId:     OPENID,
		})
	}

	// 3. 调用 Service 进行批量保存
	ctx := c.Request.Context()
	if err := service.BatchAddPlantLogs(ctx, logs, req.ImageIDs); err != nil {
		response.Fail(c, "记录失败: "+err.Error())
		return
	}

	response.Success(c, "记录成功")
}

// GetPlantLogsHandler 获取某棵植物的日志
func GetPlantLogsHandler(c *gin.Context) {
	var req GetLogsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}

	ctx := c.Request.Context()
	logs, err := service.GetPlantLogs(ctx, req.PlantID)
	if err != nil {
		response.Fail(c, "获取失败: "+err.Error())
		return
	}
	response.Success(c, logs)
}

// DeletePlantLogHandler 删除日志
func DeletePlantLogHandler(c *gin.Context) {
	var req DeleteLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}

	ctx := c.Request.Context()
	if err := service.DeletePlantLog(ctx, req.ID); err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}
	response.Success(c, "删除成功")
}

// UpdatePlantLogHandler 更新日志
func UpdatePlantLogHandler(c *gin.Context) {
	var req UpdateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}

	logTime := time.Now()
	if req.LogTime != "" {
		t, err := time.Parse("2006-01-02", req.LogTime)
		if err == nil {
			logTime = t
		}
	}

	log := &model.PlantLog{
		ActionType: req.ActionType,
		Content:    req.Content,
		LogTime:    logTime,
	}
	log.ID = req.ID

	ctx := c.Request.Context()
	if err := service.UpdatePlantLog(ctx, log, req.ImageIDs); err != nil {
		response.Fail(c, "更新失败: "+err.Error())
		return
	}

	response.Success(c, "更新成功")
}

// GetSingleLogHandler 获取单条日志
func GetSingleLogHandler(c *gin.Context) {
	var req GetSingleLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}

	ctx := c.Request.Context()
	log, err := service.GetPlantLog(ctx, req.ID)
	if err != nil {
		response.Fail(c, "获取失败: "+err.Error())
		return
	}
	response.Success(c, log)
}
