package handler

import (
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"time"

	"github.com/gin-gonic/gin"
)

type CreatePlantRequest struct {
	Name     string    `json:"name" binding:"required"`
	FamilyID uint      `json:"familyId" binding:"required"`
	TagIDs   []uint    `json:"tags"`
	Cover    string    `json:"cover"`
	Desc     string    `json:"desc"`
	Birthday time.Time `json:"birthday"`
	OpenId   string    `json:"openId"`
}

type UpdatePlantRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Cover  string `json:"cover"`
	TagIDs []uint `json:"tagIds"` // 如果不传 nil，传空数组 [] 代表清空标签
}
type GetPlantsRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

type plantRequest struct {
	ID uint `json:"id" binding:"required"`
}

func CreatePlantHandler(c *gin.Context) {
	var req CreatePlantRequest
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误"+err.Error())
		return
	}
	birthday := req.Birthday
	if birthday.IsZero() {
		birthday = time.Now()
	}
	plant := &model.Plant{
		Name:     req.Name,
		FamilyID: req.FamilyID,
		Cover:    req.Cover,
		Desc:     req.Desc,
		Birthday: birthday,
		OpenId:   OPENID,
	}
	if err := service.AddPlant(plant, req.TagIDs); err != nil {
		response.Fail(c, "新建失败:"+err.Error())
		return
	}
	response.Success(c, plant)
}

func UpdatePlantHandler(c *gin.Context) {
	var req UpdatePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误"+err.Error())
		return
	}
	updateData := make(map[string]interface{})
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Desc != "" {
		updateData["name"] = req.Desc
	}
	if req.Cover != "" {
		updateData["name"] = req.Cover
	}
	if err := service.UpdatePlant(uint(req.ID), updateData, req.TagIDs); err != nil {
		response.Fail(c, "更新失败:"+err.Error())
		return
	}
	response.Success(c, "更新成功")
}

func GetPlantsHandler(c *gin.Context) {
	var req GetPlantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误"+err.Error())
		return
	}
	plants, err := service.GetPlants(req.FamilyID)
	if err != nil {
		response.Fail(c, "获取失败:"+err.Error())
		return
	}
	response.Success(c, plants)
}

func DeletePlantHandler(c *gin.Context) {
	var req plantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误"+err.Error())
		return
	}
	if err := service.DeletePlant(req.ID); err != nil {
		response.Fail(c, "删除失败:"+err.Error())
		return
	}
	response.Success(c, "删除成功")
}

func GetPlantHandler(c *gin.Context) {
	var req plantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误"+err.Error())
		return
	}
	plant, err := service.GetPlant(req.ID)
	if err != nil {
		response.Fail(c, "获取失败:"+err.Error())
		return
	}
	response.Success(c, plant)
}
