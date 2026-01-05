package dao

import (
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

// 创建用户
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// 获取用户列表
func GetUserList(OPENID string) ([]model.User, error) {
	var user []model.User
	err := db.DB.First(&user, OPENID).Error
	return user, err
}

// 更新用户列表
func UpdateUser(user *model.User) error {
	return db.DB.Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(id int) error {
	return db.DB.Delete(&model.User{}, id).Error
}
