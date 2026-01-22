package handler

import (
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

type GetTagListRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

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

type SortTagRequest struct {
	TagIDs []uint `json:"tagIds" binding:"required"`
}

func GetTagListHandler(c *gin.Context) {
	var req GetTagListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	ctx := c.Request.Context()
	tags, err := service.GetFamilyTag(ctx, req.FamilyID)
	if err != nil {
		response.Fail(c, "获取失败:"+err.Error())
		return
	}
	response.Success(c, tags)
}

func CreateTagHandler(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误:"+err.Error())
		return
	}
	ctx := c.Request.Context()
	tag, err := service.AddTag(ctx, req.Name, req.FamilyID)
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
	ctx := c.Request.Context()
	if err := service.DeleteTag(ctx, req.ID); err != nil {
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
	ctx := c.Request.Context()
	err := service.UpdateTag(ctx, req.ID, req.Name)
	if err != nil {
		response.FailWithCode(c, 500, "更新失败"+err.Error())
		return
	}
	response.Success(c, "更新成功")
}

func SortTagHandler(c *gin.Context) {
	var req SortTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}
	ctx := c.Request.Context()
	err := service.SortTag(ctx, req.TagIDs)
	if err != nil {
		response.Fail(c, "排序失败:"+err.Error())
		return
	}
	response.Success(c, "排序保存成功")
}
