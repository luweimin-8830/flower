package dao

import (
	"context" // 假设这是你的全局db包
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

// GetByHash 根据哈希查找图片
func GetImageByHash(ctx context.Context, hash string) (*model.Image, error) {
	var img model.Image
	err := db.DB.WithContext(ctx).Where("hash = ?", hash).First(&img).Error
	return &img, err
}

// Create 创建图片记录
func CreateImage(ctx context.Context, img *model.Image) error {
	return db.DB.WithContext(ctx).Where(model.Image{Hash: img.Hash}).FirstOrCreate(img).Error
}
