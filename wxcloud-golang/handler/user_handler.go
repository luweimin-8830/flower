package handler

import (
	"fmt"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

type SortFamilyRequest struct {
	FamilyIDs []uint `json:"familyIds" binding:"required"`
}

type SwitchFamilyRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

type DeleteFamilyRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

type UpdateFamilyRequest struct {
	FamilyID uint   `json:"familyId" binding:"required"`
	Name     string `json:"name" binding:"required,min=1,max=50"`
}

type CreateFamilyRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

type JoinFamilyRequest struct {
	FamilyID uint `json:"familyId" binding:"required"`
}

type UpdateUserRequest struct {
	RemindTime string `json:"remindTime" binding:"required"`
}

type AddRemindRequest struct {
	FamilyID   uint   `json:"familyId"`
	PlantID    uint   `json:"plantId"`
	RemindTime string `json:"remindTime" binding:"required"` // 格式: 2026-03-19 08:00
	Content    string `json:"content" binding:"required"`
	ActionType string `json:"actionType"`
}

func IndexHandler(c *gin.Context) {
	response.SuccessMsg(c, "Hello Succulent")
}

func CreateUserHandler(c *gin.Context) {
	OPENID := c.GetHeader("X-WX-OPENID")
	if OPENID == "" {
		fmt.Println("⚠️ 警告: 未获取到 OpenID")
		response.FailWithCode(c, 401, "未获取到OpenId")
		return
	}

	ctx := c.Request.Context()
	user, err := service.AddUser(ctx, OPENID)
	if err != nil {
		response.FailWithCode(c, 500, "创建用户失败:"+err.Error())
		return
	}
	response.Success(c, user)
}

func UserLoginHandler(c *gin.Context) {
	OPENID := c.GetHeader("X-WX-OPENID")
	if OPENID == "" {
		fmt.Println("⚠️ 警告: 未获取到 OpenID")
		response.FailWithCode(c, 401, "未获取到OpenId")
		return
	}
	ctx := c.Request.Context()
	user, family, err := service.Login(ctx, OPENID)
	if err != nil {
		response.FailWithCode(c, 500, "登录失败:"+err.Error())
		return
	}
	response.Success(c, gin.H{
		"user":   user,
		"family": family,
	})
}

func SortFamilyHandler(c *gin.Context) {
	var req SortFamilyRequest
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}
	ctx := c.Request.Context()
	err := service.UpdateFamilySort(ctx, req.FamilyIDs, OPENID)
	if err != nil {
		response.Fail(c, "排序失败:"+err.Error())
		return
	}
	response.Success(c, "排序保存成功")
}

func SwitchFamilyHandler(c *gin.Context) {
	var req SwitchFamilyRequest
	OPENID := c.GetHeader("X-WX-OPENID")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}
	ctx := c.Request.Context()
	err := service.SwitchFamily(ctx, OPENID, req.FamilyID)
	if err != nil {
		response.Fail(c, "切换失败:"+err.Error())
		return
	}
	response.Success(c, "切换成功")
}

func DeleteFamilyHandler(c *gin.Context) {
	var req DeleteFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}
	ctx := c.Request.Context()
	err := service.DeleteFamily(ctx, req.FamilyID)
	if err != nil {
		response.Fail(c, "删除失败:"+err.Error())
		return
	}
	response.Success(c, "删除成功")
}

func GetFamilyHandler(c *gin.Context) {
	OPENID := c.GetHeader("X-WX-OPENID")
	ctx := c.Request.Context()
	family, err := service.GetFamilyList(ctx, OPENID)
	if err != nil {
		response.Fail(c, "获取失败:"+err.Error())
		return
	}
	response.Success(c, family)
}

func UpdateFamilyHandler(c *gin.Context) {
	var req UpdateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误")
		return
	}
	OPENID := c.GetHeader("X-WX-OPENID")
	ctx := c.Request.Context()
	err := service.UpdateFamily(ctx, OPENID, req.FamilyID, req.Name)
	if err != nil {
		response.Fail(c, "更新失败:"+err.Error())
		return
	}
	response.Success(c, "更新成功")

}

func CreateFamilyHandler(c *gin.Context) {
	// 1. 参数校验
	var req CreateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 400, "家庭名称不能为空")
		return
	}
	OPENID := c.GetHeader("X-WX-OPENID")
	family, err := service.CreateFamily(c.Request.Context(), OPENID, req.Name)
	if err != nil {
		response.Fail(c, "创建失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"id":   family.ID,
		"name": family.Name,
	})
}

func JoinFamilyHandler(c *gin.Context) {
	var req JoinFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 400, "家庭ID不能为空")
		return
	}
	OPENID := c.GetHeader("X-WX-OPENID")
	ctx := c.Request.Context()
	err := service.JoinFamily(ctx, OPENID, req.FamilyID)
	if err != nil {
		response.Fail(c, "加入失败: "+err.Error())
		return
	}
	response.Success(c, "加入成功")
}

func UpdateUserHandler(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 400, "参数错误")
		return
	}
	OPENID := c.GetHeader("X-WX-OPENID")
	ctx := c.Request.Context()
	err := service.UpdateUserRemindTime(ctx, OPENID, req.RemindTime)
	if err != nil {
		response.Fail(c, "更新失败: "+err.Error())
		return
	}
	response.Success(c, "更新成功")
}

func AddRemindHandler(c *gin.Context) {
	var req AddRemindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 400, "参数错误")
		return
	}

	OPENID := c.GetHeader("X-WX-OPENID")
	ctx := c.Request.Context()

	// 解析时间
	loc, _ := time.LoadLocation("Local")
	remindTime, err := time.ParseInLocation("2006-01-02 15:04", req.RemindTime, loc)
	if err != nil {
		// 尝试带秒的格式
		remindTime, err = time.ParseInLocation("2006-01-02 15:04:05", req.RemindTime, loc)
		if err != nil {
			response.Fail(c, "时间格式错误，请使用 YYYY-MM-DD HH:mm")
			return
		}
	}

	err = service.AddRemindTask(ctx, OPENID, req.FamilyID, req.PlantID, remindTime, req.Content, req.ActionType)
	if err != nil {
		response.Fail(c, "预约失败: "+err.Error())
		return
	}
	response.Success(c, "预约成功")
}
