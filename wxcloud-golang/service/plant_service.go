package service

import (
	"fmt"
	"time"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

var plantDao = new(dao.PlantDao)

func GetPlantList(req model.PlantListReq, openId string) ([]model.Plant, int64, error) {
	//处理业务逻辑

	//调用dao
	return plantDao.GetList(req, openId)
}

func AddPlant(req model.PlantAddReq, openId string) error {
	var birthday time.Time
	var err error
	if req.Birthday != "" {
		birthday, err = time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			fmt.Printf("传入日期格式错误 %+v\n", err)
			return err
		}
	} else {
		birthday = time.Now()
	}
	plant := model.Plant{
		Name:     req.Name,
		Cover:    req.Cover,
		Desc:     req.Desc,
		Labels:   req.Labels,
		Tags:     req.Tags,
		OpenId:   openId,
		Birthday: birthday,
	}

	return plantDao.Create(&plant)
}
