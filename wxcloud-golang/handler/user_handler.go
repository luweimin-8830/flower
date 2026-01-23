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
