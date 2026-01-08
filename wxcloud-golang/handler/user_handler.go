package handler

import (
	"fmt"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

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

	user, err := service.AddUser(OPENID)
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
	user, family, err := service.Login(OPENID)
	if err != nil {
		response.FailWithCode(c, 500, "登录失败:"+err.Error())
		return
	}
	response.Success(c, gin.H{
		"user":   user,
		"family": family,
	})
}
