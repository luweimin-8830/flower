package handler

import (
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
	"wxcloud-golang/service"

	"github.com/gin-gonic/gin"
)

// CheckRequest 检查请求
type CheckRequest struct {
	Hash string `json:"hash" binding:"required"`
}

// SaveImageRequest 保存请求
type SaveImageRequest struct {
	URL    string `json:"url" binding:"required"` // 云存储 FileID
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int64  `json:"size"`
	Hash   string `json:"hash" binding:"required"`
}

// CheckImageHandler 1. 检查图片是否已存在
func CheckImageHandler(c *gin.Context) {
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	img, needUpload, err := service.CheckImageExist(c.Request.Context(), req.Hash)
	if err != nil {
		response.Fail(c, "查询失败")
		return
	}

	if !needUpload {
		// 图片已存在，直接返回 ID 和 URL
		// 前端拿到这个 ID 后，直接去调创建日志接口，跳过上传步骤
		response.Success(c, gin.H{
			"uploadRequired": false,
			"id":             img.ID,
			"url":            img.URL,
		})
		return
	}

	// 图片不存在，告诉前端去上传
	response.Success(c, gin.H{
		"uploadRequired": true,
	})
}

// SaveImageHandler 2. 上传完成后，保存元数据
func SaveImageHandler(c *gin.Context) {
	var req SaveImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	img := &model.Image{
		URL:    req.URL,
		Width:  req.Width,
		Height: req.Height,
		Hash:   req.Hash,
		Size:   req.Size, // 赋值 size
	}

	if err := service.SaveImageMetadata(c.Request.Context(), img); err != nil {
		response.Fail(c, "保存失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":  img.ID,
		"url": img.URL,
	})
}
