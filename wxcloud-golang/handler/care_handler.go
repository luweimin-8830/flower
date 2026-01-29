package handler

import (
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

type GetCareRequest struct {
	FamilyID uint `json:"familyId"`
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
		response.Fail(c, "参数错误: "+err.Error())
		return
	}
	ctx := c.Request.Context()
	
	// 如果前端传 0，尝试获取该用户的默认家庭
	familyID := req.FamilyID
	if familyID == 0 {
		openID := c.GetHeader("X-WX-OPENID")
		user, err := service.GetUserWithFamilies(ctx, openID)
		if err == nil && len(user.Families) > 0 {
			familyID = user.Families[0].ID
		}
	}

	if familyID == 0 {
		response.Fail(c, "家庭ID无效")
		return
	}

	cares, err := service.GetFamilyCareActions(ctx, familyID)
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
