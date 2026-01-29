package handler

import (
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

type GetCareRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

type AddCareRequest struct {
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon"`
	Color    string `json:"color"`
	Type     string `json:"type" binding:"required"`
	FamilyID uint   `json:"familyId" binding:"required"`
}

type UpdateCareRequest struct {
	ID    uint   `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Type  string `json:"type" binding:"required"`
}

type DeleteCareRequest struct {
	ID uint `json:"id" binding:"required"`
}

type SortCareRequest struct {
	CareIDs []uint `json:"careIds" binding:"required"`
}

func GetCareListHandler(c *gin.Context) {
	var req GetCareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	ctx := c.Request.Context()
	cares, err := service.GetFamilyCareActions(ctx, req.FamilyID)
	if err != nil {
		response.Fail(c, "获取失败: "+err.Error())
		return
	}
	response.Success(c, cares)
}

func AddCareHandler(c *gin.Context) {
	var req AddCareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	ctx := c.Request.Context()
	care := &model.CareAction{
		Name:     req.Name,
		Icon:     req.Icon,
		Color:    req.Color,
		Type:     req.Type,
		FamilyID: req.FamilyID,
	}
	if err := service.AddCareAction(ctx, care); err != nil {
		response.Fail(c, "添加失败: "+err.Error())
		return
	}
	response.Success(c, care)
}

func UpdateCareHandler(c *gin.Context) {
	var req UpdateCareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	ctx := c.Request.Context()
	care := &model.CareAction{
		Name:  req.Name,
		Icon:  req.Icon,
		Color: req.Color,
		Type:  req.Type,
	}
	care.ID = req.ID
	if err := service.UpdateCareAction(ctx, care); err != nil {
		response.Fail(c, "更新失败: "+err.Error())
		return
	}
	response.Success(c, "更新成功")
}

func DeleteCareHandler(c *gin.Context) {
	var req DeleteCareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	ctx := c.Request.Context()
	if err := service.DeleteCareAction(ctx, req.ID); err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}
	response.Success(c, "删除成功")
}

func SortCareHandler(c *gin.Context) {
	var req SortCareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	ctx := c.Request.Context()
	if err := service.SortCareActions(ctx, req.CareIDs); err != nil {
		response.Fail(c, "排序失败: "+err.Error())
		return
	}
	response.Success(c, "排序成功")
}
