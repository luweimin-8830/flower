package service

import (
	"errors"
	"fmt"
	"time"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// AddUser 新增业务
func AddUser(openId string) (*model.User, error) {
	user := &model.User{
		OPENID:     openId,
		LastDateAT: time.Now(),
		CreatedAT:  time.Now(),
	}
	err := dao.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

// 小程序登录逻辑
func Login(openId string) (*model.User, []model.Family, error) {
	var user *model.User
	var err error
	user, err = dao.GetUserByOpenID(openId)
	if err != nil {
		//未找到,新用户
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("新用户: %+v\n", err)
			user, err = AddUser(openId)
			if err != nil {
				return nil, nil, fmt.Errorf("创建用户失败: %v", err)
			}
		} else {
			return nil, nil, err
		}
	} else {
		if updateErr := dao.UpdateUserLastLogin(user.ID); updateErr != nil {
			fmt.Printf("更新登录时间失败: %v\n", updateErr)
		}
		user.LastDateAT = time.Now()
	}

	_, err = dao.GetFamilyByOpenId(openId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newFamily := &model.Family{
			Name:        "我的花园",
			OwnerOpenId: openId,
		}
		if err := dao.CreateFamily(newFamily); err != nil {
			return nil, nil, fmt.Errorf("创建家庭失败: %v", err)
		}
		newMember := &model.FamilyMember{
			FamilyID: newFamily.ID,
			OpenID:   openId,
			Role:     "owner",
		}
		if err := dao.CreateFamilyMember(newMember); err != nil {
			return nil, nil, fmt.Errorf("创建家庭成员失败: %v", err)
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("查询家庭失败: %v", err)
	}
	family, err := dao.GetFamilyList(openId)
	if err != nil {
		return nil, nil, fmt.Errorf("获取或创建家庭失败: %v", err)
	}
	return user, family, nil
}
