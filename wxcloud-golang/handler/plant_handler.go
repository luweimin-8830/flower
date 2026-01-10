package handler

import (
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"time"

	"github.com/gin-gonic/gin"
)

type CreatePlantRequest struct {
	Name     string `json:"name" binding:"required"`
	FamilyID uint   `json:"familyId" binding:"required"`
	Tags     []struct {
		ID uint `json:"id"`
	} `json:"tags"`
	Cover    string    `json:"cover"`
	Desc     string    `json:"desc"`
	Birthday time.Time `json:"birthday"`
	OpenId   string    `json:"openId"`
}

// {
// 	"familyId":1,
// 	"name":"植物2",
// 	"tags":[{"id":1},{"id":2}]
//   }

type UpdatePlantRequest struct {
	ID       uint      `json:"id" binding:"required"`
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Cover    string    `json:"cover"`
	Birthday time.Time `json:"birthday"`
	Tags     []struct {
		ID uint `json:"id"`
	} `json:"tags"` // 如果不传 nil，传空数组 [] 代表清空标签
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

	var tagIDs []uint
	for _, t := range req.Tags {
		tagIDs = append(tagIDs, t.ID)
	}

	plant := &model.Plant{
		Name:     req.Name,
		FamilyID: req.FamilyID,
		Cover:    req.Cover,
		Desc:     req.Desc,
		Birthday: birthday,
		OpenId:   OPENID,
	}
	if err := service.AddPlant(plant, tagIDs); err != nil {
		response.Fail(c, "新建失败:"+err.Error())
		return
	}

	fullPlant, err := service.GetPlant(plant.ID)
	if err != nil {
		response.Success(c, plant)
	}

	response.Success(c, fullPlant)
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
		updateData["desc"] = req.Desc
	}
	if req.Cover != "" {
		updateData["cover"] = req.Cover
	}
	if !req.Birthday.IsZero() {
		updateData["birthday"] = req.Birthday
	}

	var tagIDs []uint
	if req.Tags != nil {
		tagIDs = make([]uint, 0)
		for _, t := range req.Tags {
			tagIDs = append(tagIDs, t.ID)
		}
	}

	if err := service.UpdatePlant(uint(req.ID), updateData, tagIDs); err != nil {
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
