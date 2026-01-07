package handler

import (
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

func PlantListHandler(c *gin.Context) {
	var req model.PlantListReq
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	if OPENID == "" {
		response.FailWithCode(c, 401, "未获取到OpenId")
		return
	}

	list, total, err := service.GetPlantList(req, OPENID)

	if err != nil {
		response.Fail(c, "获取列表失败:"+err.Error())
		return
	}

	response.Success(c, gin.H{
		"data":  list,
		"total": total,
	})

}

func AddPlantHandler(c *gin.Context) {
	var req model.PlantAddReq
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	if OPENID == "" {
		response.FailWithCode(c, 401, "未获取到OpenId")
		return
	}
	err := service.AddPlant(req, OPENID)
	if err != nil {
		response.Fail(c, "写入失败"+err.Error())
		return
	}

	response.Success(c, "写入成功")
}

func DeletePlantHandler(c *gin.Context) {
	var req model.PlantDeleteReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	err := service.DeletePlant(req.ID)
	if err != nil {
		response.Fail(c, "删除失败"+err.Error())
	}
	response.Success(c, "删除成功")
}

func UpdatePlantHandler(c *gin.Context) {
	var req model.PlantUpdateReq
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	if OPENID == "" {
		response.FailWithCode(c, 401, "未获取到OpenId")
		return
	}

	if err := service.UpdatePlant(req, OPENID); err != nil {
		response.Fail(c, "更新失败"+err.Error())
		return
	}
	response.Success(c, "更新成功")
}
