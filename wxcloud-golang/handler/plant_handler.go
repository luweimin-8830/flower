package handler

import (
	"wxcloud-golang/response"
	"wxcloud-golang/service"
	"wxcloud-golang/db/model"


	"github.com/gin-gonic/gin"
)

func PlantListHandler(c *gin.Context) {
	var req model.PlantListReq
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithCode(c,401,"参数错误")
		return
	}
	if OPENID == "" {
		response.FailWithCode(c, 401, "未获取到OpenId")
		return
	}

	list, total,err := service.GetPlantList(req,OPENID)

	if err != nil {
		response.Fail(c,"获取列表失败")
		return
	}

	response.Success(c,gin.H{
		"data": list,
		"total": total,
	})
	
}