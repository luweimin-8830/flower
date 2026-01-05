package dao

import (
	"time"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
)

// 创建用户
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// 获取用户列表
func GetUserByOpenID(OPENID string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("open_id = ?", OPENID).First(&user).Error
	return &user, err
}

// 更新用户列表
func UpdateUser(user *model.User) error {
	return db.DB.Save(user).Error
}

// 更新最后登录时间
func UpdateUserLastLogin(userID uint) error {
	return db.DB.Model(&model.User{}).Where("id = ?", userID).Update("last_date_at", time.Now()).Error
}

// DeleteUser 删除用户
func DeleteUser(id int) error {
	return db.DB.Delete(&model.User{}, id).Error
}
