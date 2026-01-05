package service

import (
	"time"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"
)

// AddUser 新增业务
func AddUser(openId string) (*model.User, error) {
	user := &model.User{
		OPENID:     openId,
		LastDateAT: time.Now(),
		CreatedAT:  time.Now(),
	}
	err := dao.CreateUser(user)
	return user, err
}
