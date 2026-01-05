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
func Login(openId string) (*model.User, error) {
	user, err := dao.GetUserByOpenID(openId)
	if err != nil {
		//未找到,新用户
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("新用户: %+v\n", err)
			return AddUser(openId)
		}
		return nil, err
	}
	if updateErr := dao.UpdateUserLastLogin(user.ID); updateErr != nil {
		return nil, updateErr
	}

	user.LastDateAT = time.Now()
	return user, nil
}
