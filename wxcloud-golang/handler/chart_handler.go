package handler

import (
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
