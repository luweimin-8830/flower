package dao

import (
	"errors"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

type PlantDao struct {
}

func (d *PlantDao) GetList(req model.PlantListReq, openId string) ([]model.Plant, int64, error) {
	var plants []model.Plant
	var total int64
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.PageSize
	if size <= 0 {
		size = 10
	}
	selectSQL := `plants.*,IF(plants.open_id = ?,'owner', plant_user_relations.role) as current_user_role`
	err := db.DB.Model(&model.Plant{}).
		Joins("LEFT JOIN plant_user_relations ON plant_user_relations.plant_id = plants.id AND plant_user_relations.open_id = ?", openId).
		Select(selectSQL, openId).Where("plants.open_id = ? OR plant_user_relations.open_id = ?", openId, openId).Offset((page - 1) * size).Limit(size).
		Order("plants.updated_at DESC").Find(&plants).Error
	return plants, total, err
}

// 创建新植物
func (d *PlantDao) Create(plant *model.Plant) error {
	return db.DB.Create(plant).Error
}

// 删除植物
func (d *PlantDao) Delete(id uint) error {
	return db.DB.Where("id = ?", id).Delete(&model.Plant{}).Error
}

func (d *PlantDao) Update(id uint, openId string, updates map[string]interface{}) error {
	result := db.DB.Model(&model.Plant{}).Where("id = ? AND open_id = ?", id, openId).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("更新失败,数据不存在或无权限")
	}

	return nil
}

// func (d *PlantDao) Check(plantID uint, openId string, requiredAction string) (bool, error) {
// 	var relation model.PlantUserRelation

// 	err := db.DB.Where("plant_id = ? AND open_id = ?", plantID, openId).First(&relation).Error

// 	if err != nil {
// 		var plant model.Plant
// 		if err := db.DB.Select("open_id").First(&plant, plantID).Error; err != nil {
// 			return false, err
// 		}
// 		if plant.OpenId == openId {
// 			return true, nil
// 		}
// 		return false, nil
// 	}

// 	switch requiredAction {
// 	case "view":
// 		return true, nil
// 	case "edit":
// 		return relation.Role == "editor" || relation.Role == "admin", nil
// 	case "delete":
// 		return relation.Role == "admin", nil
// 	}

// 	return false, nil
// }
