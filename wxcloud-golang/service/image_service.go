package service

import (
	"context"
	"errors"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// CheckImageExist 检查图片是否存在
// 返回: (image对象, 是否需要上传, error)
func CheckImageExist(ctx context.Context, sha256 string) (*model.Image, bool, error) {
	img, err := dao.GetImageByHash(ctx, sha256)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到 -> 需要上传
			return nil, true, nil
		}
		// 数据库报错
		return nil, false, err
	}
	// 找到了 -> 不需要上传
	return img, false, nil
}

// SaveImageMetadata 保存图片元数据
func SaveImageMetadata(ctx context.Context, img *model.Image) error {
	return dao.CreateImage(ctx, img)
}
