package handler

import (
	"time"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"

	"github.com/gin-gonic/gin"
)

type ChartDataRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

type ChartItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func GetChartDataHandler(c *gin.Context) {
	var req ChartDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}

	ctx := c.Request.Context()

	// 1. 获取所有有标签的统计
	var tagStats []ChartItem
	err := db.DB.WithContext(ctx).Model(&model.Tag{}).
		Select("tag.name as name, COUNT(plant.id) as value").
		Joins("JOIN plant_tags ON plant_tags.tag_id = tag.id").
		Joins("JOIN plant ON plant.id = plant_tags.plant_id AND plant.deleted_at IS NULL").
		Where("tag.family_id = ?", req.FamilyID).
		Group("tag.id").
		Find(&tagStats).Error

	if err != nil {
		response.Fail(c, "获取统计数据失败")
		return
	}

	// 2. 获取该家庭下总植物数
	var totalPlants int64
	db.DB.WithContext(ctx).Model(&model.Plant{}).Where("family_id = ?", req.FamilyID).Count(&totalPlants)

	// 3. 计算有标签的植物总数（注意：一个植物可能有多个标签，这里统计的是标签关联数，但用户需求通常是“占比”）
	// 为了实现“没标签的塞入其他”，我们需要计算“没有任何标签关联的植物数量”
	var taggedPlantCount int64
	db.DB.WithContext(ctx).Model(&model.Plant{}).
		Joins("JOIN plant_tags ON plant_tags.plant_id = plant.id").
		Where("plant.family_id = ? AND plant.deleted_at IS NULL", req.FamilyID).
		Distinct("plant.id").
		Count(&taggedPlantCount)

	otherCount := totalPlants - taggedPlantCount
	if otherCount > 0 {
		tagStats = append(tagStats, ChartItem{
			Name:  "其他",
			Value: otherCount,
		})
	}

	response.Success(c, tagStats)
}

type CareChartRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
	Year     int  `json:"year" binding:"required"`
	Month    int  `json:"month" binding:"required"`
}

type CareChartResponse struct {
	Categories []string `json:"categories"`
	Series     []Series `json:"series"`
}

type Series struct {
	Name string  `json:"name"`
	Data []int64 `json:"data"`
}

func GetCareChartHandler(c *gin.Context) {
	var req CareChartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}

	ctx := c.Request.Context()

	// 1. 获取该家庭在该月份的所有养护记录统计
	// 我们按 ActionType 统计每天的次数。为了简化，我们展示该月内前 10 种最常见的操作类型
	var results []struct {
		ActionType string
		Count      int64
	}

	// 构造月份范围
	startDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	err := db.DB.WithContext(ctx).Model(&model.PlantLog{}).
		Select("action_type, COUNT(*) as count").
		Joins("JOIN plant ON plant.id = plant_log.plant_id").
		Where("plant.family_id = ? AND plant_log.log_time >= ? AND plant_log.log_time < ? AND plant_log.deleted_at IS NULL", req.FamilyID, startDate, endDate).
		Group("action_type").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error

	if err != nil {
		response.Fail(c, "获取养护统计失败")
		return
	}

	resp := CareChartResponse{
		Categories: []string{},
		Series:     []Series{{Name: "操作次数", Data: []int64{}}},
	}

	for _, res := range results {
		resp.Categories = append(resp.Categories, model.GetActionName(res.ActionType))
		resp.Series[0].Data = append(resp.Series[0].Data, res.Count)
	}

	response.Success(c, resp)
}
