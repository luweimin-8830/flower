package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wxcloud-golang/db/dao"
	"wxcloud-golang/db/model"

	"gorm.io/gorm"
)

// AddUser 新增业务
func AddUser(ctx context.Context, openId string) (*model.User, error) {
	user := &model.User{
		OPENID:     openId,
		LastDateAT: time.Now(),
		CreatedAT:  time.Now(),
	}
	err := dao.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

// 小程序登录逻辑
func Login(ctx context.Context, openId string) (*model.User, []model.Family, error) {
	var user *model.User
	var err error
	user, err = dao.GetUserByOpenID(ctx, openId)
	if err != nil {
		//未找到,新用户
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("新用户: %+v\n", err)
			user, err = AddUser(ctx, openId)
			if err != nil {
				return nil, nil, fmt.Errorf("创建用户失败: %v", err)
			}
		} else {
			return nil, nil, err
		}
	} else {
		if updateErr := dao.UpdateUserLastLogin(ctx, user.ID); updateErr != nil {
			fmt.Printf("更新登录时间失败: %v\n", updateErr)
		}
		user.LastDateAT = time.Now()
	}

	_, err = dao.GetFamilyByOpenId(ctx, openId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newFamily := &model.Family{
			Name:        "我的花园",
			OwnerOpenId: openId,
		}
		if err := dao.CreateFamily(ctx, newFamily); err != nil {
			return nil, nil, fmt.Errorf("创建家庭失败: %v", err)
		}
		newMember := &model.FamilyMember{
			FamilyID: newFamily.ID,
			OpenID:   openId,
			Role:     "owner",
		}
		if err := dao.CreateFamilyMember(ctx, newMember); err != nil {
			return nil, nil, fmt.Errorf("创建家庭成员失败: %v", err)
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("查询家庭失败: %v", err)
	}
	family, err := dao.GetFamilyList(ctx, openId)
	if err != nil {
		return nil, nil, fmt.Errorf("获取或创建家庭失败: %v", err)
	}
	return user, family, nil
}

func UpdateFamilySort(ctx context.Context, familyIDs []uint, OPENID string) error {
	return dao.UpdateFamilySortOrder(ctx, familyIDs, OPENID)
}

func SwitchFamily(ctx context.Context, openID string, familyID uint) error {
	return dao.SwitchCurrentFamily(ctx, openID, familyID)
}

func DeleteFamily(ctx context.Context, familyID uint) error {
	return dao.DeleteFamilyWithData(ctx, familyID)
}
