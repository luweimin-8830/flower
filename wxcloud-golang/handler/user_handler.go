package handler

import (
	"fmt"
	"net/http"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	OPENID := c.GetHeader("X-WX-OPENID")
	if OPENID == "" {
		fmt.Println("⚠️ 警告: 未获取到 OpenID，使用测试账号")
		c.JSON(400, gin.H{"error": "未获取到OpenId"})
		return
	}

	user, err := service.AddUser(OPENID)
	if err != nil {
		c.JSON(500, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}
