package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // 必须导入，否则无法识别 jpg
	_ "image/png"  // 必须导入，否则无法识别 png
	"io"
	"path"

	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
	"wxcloud-golang/utils"

	"gorm.io/gorm"
)

const (
	WxEnvID      = "prod-0gr2o3qpe533f1fb"                 // 你的环境ID
	WxBucketName = "7072-prod-0gr2o3qpe533f1fb-1352691102" // 你的存储桶名称
)

// UploadImage 处理图片上传的核心业务
// fileData 必须是 io.ReadSeeker，因为我们需要多次读取流（算Hash、取宽高、上传）
func UploadImage(ctx context.Context, fileData io.ReadSeeker, originalName string, familyID uint) (*model.Image, error) {
	// 1. 计算 Hash (用于秒传/查重)
	hashObj := sha256.New()
	if _, err := io.Copy(hashObj, fileData); err != nil {
		return nil, errors.New("hash计算失败")
	}
	hash := hex.EncodeToString(hashObj.Sum(nil))

	// 2. 【查重】查询数据库
	existImg, err := dao.Image.GetByHash(ctx, hash)
	if err == nil {
		return existImg, nil // 秒传成功：数据库已有，直接返回旧对象
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err // 真正的数据库错误
	}

	// --- 此时指针在末尾，必须重置 ---
	fileData.Seek(0, 0)

	// 3. 获取图片宽高 & 格式校验
	config, format, err := image.DecodeConfig(fileData)
	if err != nil {
		return nil, errors.New("无法识别图片格式或文件损坏")
	}

	// --- 再次重置指针 ---
	fileData.Seek(0, 0)

	// 4. 获取文件大小 (通过移动指针到末尾)
	size, _ := fileData.Seek(0, io.SeekEnd)
	fileData.Seek(0, 0) // 记得移回去！

	// 5. 构造存储路径 (按 FamilyID 分文件夹)
	ext := path.Ext(originalName)
	if ext == "" {
		ext = "." + format // 如果原文件名没后缀，用识别出的格式补全
	}

	// 逻辑：如果有 familyID，存入 families/101/hash.jpg，否则存入 common/hash.jpg
	var objectKey string
	if familyID > 0 {
		objectKey = fmt.Sprintf("families/%d/%s%s", familyID, hash, ext)
	} else {
		objectKey = fmt.Sprintf("common/%s%s", hash, ext)
	}

	// 6. 上传到 COS
	// utils.UploadToCOS 内部应该接收 io.Reader
	_, err = utils.UploadToCOS(fileData, objectKey)
	if err != nil {
		fmt.Printf("COS Upload Error: %v\n", err)
		return nil, err
	}

	url := fmt.Sprintf("cloud\\://%s.%s/%s", WxEnvID, WxBucketName, objectKey)

	// 7. 入库
	newImg := &model.Image{
		URL:      url,
		Hash:     hash,
		Width:    config.Width,
		Height:   config.Height,
		Size:     size, // 使用计算出的 size
		MimeType: "image/" + format,
	}

	// 如果你的 Image 模型里有 FamilyID 字段，记得在这里赋值
	// newImg.FamilyID = familyID

	if err := dao.Image.Create(ctx, newImg); err != nil {
		return nil, err
	}

	return newImg, nil
}
