package handler

import (
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

type CreateTagRequest struct {
	Name     string `json:"name" binding:"required"`
	FamilyID uint   `json:"familyId" binding:"required"`
}

type DeleteTagRequest struct {
	ID uint `json:"id" binding:"required"`
}

type UpdateTagRequest struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func CreateTagHandler(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	tag, err := service.AddTag(req.Name, req.FamilyID)
	if err != nil {
		response.FailWithCode(c, 500, "创建失败"+err.Error())
		return
	}
	response.Success(c, tag)
}

func DeleteTagHandler(c *gin.Context) {
	var req DeleteTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	if err := service.DeleteTag(req.ID); err != nil {
		response.FailWithCode(c, 500, "删除失败"+err.Error())
		return
	}
	response.Success(c, nil)
}

func UpdateTagHandler(c *gin.Context) {
	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	err := service.UpdateTag(req.ID, req.Name)
	if err != nil {
		response.FailWithCode(c, 500, "更新失败"+err.Error())
		return
	}
	response.Success(c, "更新成功")
}
