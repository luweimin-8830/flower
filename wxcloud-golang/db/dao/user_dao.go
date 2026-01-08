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

// 创建家庭
func CreateFamily(family *model.Family) error {
	return db.DB.Create(family).Error
}

// 创建家庭成员
func CreateFamilyMember(familyMember *model.FamilyMember) error {
	return db.DB.Create(familyMember).Error
}

// 查询家庭列表
func GetFamilyList(OPENID string) ([]model.Family, error) {
	var family []model.Family

	err := db.DB.Model(&model.Family{}).Joins("INNER JOIN family_member ON family_member.family_id = family.id").
		Where("family_member.open_id = ?", OPENID).Find(&family).Error

	return family, err
}

// 查询是否存在家庭
func GetFamilyByOpenId(OPENID string) (*model.Family, error) {
	var family model.Family

	err := db.DB.Model(&model.Family{}).Joins("INNER JOIN family_member ON family_member.family_id = family.id").
		Where("family_member.open_id = ?", OPENID).First(&family).Error

	return &family, err
}
