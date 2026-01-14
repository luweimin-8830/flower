package db

import (
	"context"

	"gorm.io/gorm"
)

// Conn 返回绑定了指定 Context 的 *gorm.DB，确保链路能够在数据库层传递。
func Conn(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return DB
	}
	return DB.WithContext(ctx)
}
