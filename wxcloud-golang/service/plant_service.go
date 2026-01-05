package service

import (
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

var plantDao = new(dao.PlantDao)

func GetPlantList(req model.PlantListReq,openId string)([]model.Plant, int64, error) {
	//处理业务逻辑

	//调用dao
	return plantDao.GetList(req,openId)
}