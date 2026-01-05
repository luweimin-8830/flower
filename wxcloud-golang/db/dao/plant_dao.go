package dao

import (
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

type PlantDao struct {
}

func (d *PlantDao) GetList(req model.PlantListReq) (plants []model.Plant, total int64, err error) {
	query := db.DB.Model(&model.Plant{})

	if req.Labels != "" {
		query = query.Where("labels = ?", req.Labels)
	}
}
