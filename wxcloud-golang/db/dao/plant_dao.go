package dao

import (
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

type PlantDao struct {
}

func (d *PlantDao) GetList(req model.PlantListReq, openId string) (plants []model.Plant, total int64, err error) {
	query := db.DB.Model(&model.Plant{})

	if openId != "" {
		query = query.Where("open_id = ?", openId)
	}

	if req.Labels != "" {
		query = query.Where("labels = ?", req.Labels)
	}

	if err = query.Count(&total).Error; err != nil {
		return
	}

	size := req.PageSize
	if size <= 0 {
		size = 10
	}
	offset := (req.Page - 1) * size

	err = query.Order("id desc").Offset(offset).Limit(size).Find(&plants).Error
	return
}

// 创建新植物
func (d *PlantDao) Create(plant *model.Plant) error {
	return db.DB.Create(plant).Error
}
