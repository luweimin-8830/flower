package handler

import (
	"bytes"
	"encoding/base64"
	"strings"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

type UploadJSONRequest struct {
	FamilyID    uint   `json:"familyId"`
	FileName    string `json:"fileName"` // 需要文件名来获取后缀
	ImageBase64 string `json:"image"`    // Base64 字符串
}

func UploadHandler(c *gin.Context) {
	var req UploadJSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, 401, "参数错误"+err.Error())
		return
	}

	if req.ImageBase64 == "" {
		response.FailWithCode(c, 401, "image不能为空")
	}

	b64Data := req.ImageBase64
	if idx := strings.Index(b64Data, ","); idx > -1 {
		b64Data = b64Data[idx+1:]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		response.FailWithCode(c, 500, "图片Base64解码失败")
		return
	}

	fileReader := bytes.NewReader(imgBytes)

	imgModel, err := service.UploadImage(c.Request.Context(), fileReader, req.FileName, req.FamilyID)
	if err != nil {
		response.Fail(c, "上传失败"+err.Error())
		return
	}

	// 3. 返回成功
	response.Success(c, imgModel)
}
