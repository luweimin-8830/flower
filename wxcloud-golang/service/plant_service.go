package service

import (
	"errors"
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
	fmt.Printf("传入参数 %+v\n", req)
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
		FamilyID: req.FamilyID,
		OpenId:   openId,
		Birthday: birthday,
	}

	return plantDao.Create(&plant)
}

func DeletePlant(id uint) error {
	return plantDao.Delete(id)
}

func UpdatePlant(req model.PlantUpdateReq, openId string) error {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Cover != nil {
		updates["cover"] = *req.Cover
	}
	if req.Desc != nil {
		updates["desc"] = *req.Desc
	}
	if req.Labels != nil {
		updates["labels"] = req.Labels
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Birthday != nil {
		if *req.Birthday == "" {
			// 如果传了空字符串，意为清空生日
			updates["birthday"] = nil
		} else {
			t, err := time.Parse("2006-01-02", *req.Birthday)
			if err != nil {
				return errors.New("日期格式错误")
			}
			updates["birthday"] = t
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return plantDao.Update(req.ID, openId, updates)
}
