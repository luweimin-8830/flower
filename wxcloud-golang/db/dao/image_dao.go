package dao

import (
	"context" // 假设这是你的全局db包
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

type imageDao struct{}

var Image = new(imageDao) // 单例模式

// GetByHash 根据哈希查找图片
func (d *imageDao) GetByHash(ctx context.Context, hash string) (*model.Image, error) {
	var img model.Image
	// 如果 execWithSpan 支持查询回调：
	err := execWithSpan(ctx, "SELECT", "image", func(conn *gorm.DB) error {
		return conn.Where("hash = ?", hash).First(&img).Error
	})
	// 如果 execWithSpan 不好用，直接用 db.DB 也可以：
	// err := db.DB.WithContext(ctx).Where("hash = ?", hash).First(&img).Error
	return &img, err
}

// Create 创建图片记录
func (d *imageDao) Create(ctx context.Context, img *model.Image) error {
	return execWithSpan(ctx, "INSERT", "image", func(conn *gorm.DB) error {
		return conn.Create(img).Error
	})
}
