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
	// 1. 获取或创建用户
	user, err = dao.GetUserByOpenID(ctx, openId)
	if err != nil {
		// 未找到, 新用户
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
		// 老用户，更新登录时间
		if updateErr := dao.UpdateUserLastLogin(ctx, user.ID); updateErr != nil {
			fmt.Printf("更新登录时间失败: %v\n", updateErr)
		}
		user.LastDateAT = time.Now()
	}
	families, err := dao.GetFamilyList(ctx, openId)
	if err != nil {
		return nil, nil, fmt.Errorf("查询家庭列表失败: %v", err)
	}
	if len(families) == 0 {
		newFamily := &model.Family{
			Name:        "我的花园",
			OwnerOpenId: openId,
		}
		// 2. 在 family_member 表添加你为 owner
		// 3. 更新 user 表的 current_family_id
		if err := dao.CreateFamily(ctx, newFamily); err != nil {
			return nil, nil, fmt.Errorf("创建默认家庭失败: %v", err)
		}
		// 初始化默认养护项
		_ = dao.CreateDefaultCareActions(ctx, newFamily.ID)
		families, err = dao.GetFamilyList(ctx, openId)
		if err != nil {
			return nil, nil, fmt.Errorf("重新获取家庭列表失败: %v", err)
		}
		user.CurrentFamilyID = &newFamily.ID
	}
	return user, families, nil
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

func GetFamilyList(ctx context.Context, OPENID string) ([]model.Family, error) {
	return dao.GetFamilyList(ctx, OPENID)
}

func UpdateFamily(ctx context.Context, openID string, familyID uint, newName string) error {
	return dao.UpdateFamilyName(ctx, openID, familyID, newName)
}

func CreateFamily(ctx context.Context, openID string, name string) (*model.Family, error) {
	family := &model.Family{
		Name:        name,
		OwnerOpenId: openID,
	}
	err := dao.CreateFamily(ctx, family)
	if err != nil {
		return nil, err
	}
	// 初始化默认养护项
	_ = dao.CreateDefaultCareActions(ctx, family.ID)
	return family, nil
}

func JoinFamily(ctx context.Context, openID string, familyID uint) error {
	// 1. 检查是否已经是成员
	families, err := dao.GetFamilyList(ctx, openID)
	if err != nil {
		return err
	}
	for _, f := range families {
		if f.ID == familyID {
			return errors.New("你已经是该家庭成员了")
		}
	}

	// 2. 加入家庭
	member := &model.FamilyMember{
		FamilyID:  familyID,
		OpenID:    openID,
		Role:      "member",
		SortOrder: len(families),
	}
	return dao.CreateFamilyMember(ctx, member)
}
