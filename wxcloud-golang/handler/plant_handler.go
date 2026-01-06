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
