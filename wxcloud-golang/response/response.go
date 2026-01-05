package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构体
type Response struct {
	Code int         `json:"code"`           // 业务状态码：0成功，其他失败
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 数据，omitempty表示如果为空则不返回该字段
}

// Result 基础封装方法
func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	Result(c, 0, "success", data)
}

// SuccessMsg 成功响应（仅提示，无数据）
func SuccessMsg(c *gin.Context, msg string) {
	Result(c, 0, msg, nil)
}

// Fail 失败响应
func Fail(c *gin.Context, msg string) {
	// 这里默认业务错误码为 -1，你也可以根据需要传入具体错误码
	Result(c, -1, msg, nil)
}

// FailWithCode 失败响应（自定义错误码）
func FailWithCode(c *gin.Context, code int, msg string) {
	Result(c, code, msg, nil)
}
