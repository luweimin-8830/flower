package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

type UserWithFamilies struct {
	*model.User
	Families []model.Family `json:"families"`
}

func GetUserWithFamilies(ctx context.Context, openID string) (*UserWithFamilies, error) {
	user, err := dao.GetUserByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}
	families, err := dao.GetFamilyList(ctx, openID)
	if err != nil {
		return nil, err
	}
	return &UserWithFamilies{
		User:     user,
		Families: families,
	}, nil
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

func UpdateUserRemindTime(ctx context.Context, openID string, remindTime string) error {
	user, err := dao.GetUserByOpenID(ctx, openID)
	if err != nil {
		return err
	}
	user.RemindTime = remindTime
	return dao.UpdateUser(ctx, user)
}

func AddRemindTask(ctx context.Context, openID string, familyID uint, plantID uint, remindTime time.Time, content string, actionType string) error {
	task := &model.RemindTask{
		OpenID:     openID,
		FamilyID:   familyID,
		PlantID:    plantID,
		RemindTime: remindTime,
		Content:    content,
		ActionType: actionType,
		Status:     0,
	}
	return dao.GetDB(ctx).Create(task).Error
}

type SubscribeMessageRequest struct {
	ToUser           string            `json:"touser"`
	TemplateID       string            `json:"template_id"`
	Page             string            `json:"page"`
	Data             map[string]Value  `json:"data"`
	MiniprogramState string            `json:"miniprogram_state"`
	Lang             string            `json:"lang"`
}

type Value struct {
	Value string `json:"value"`
}

func SendSubscribeMessage(ctx context.Context, task *model.RemindTask) error {
	templateID := "jtmMCRDxoFP3AEDRAi0yNGo9PXI_3BGyb7bcqdlSJk4"
	
	plantName := "植物养护提醒"
	if task.PlantID > 0 {
		var plant model.Plant
		err := dao.GetDB(ctx).First(&plant, task.PlantID).Error
		if err == nil && plant.Name != "" {
			plantName = plant.Name
		}
	}

	// 组装订阅消息内容
	data := map[string]Value{
		"thing8": {Value: plantName}, // 作物名称
		"thing4": {Value: "需要关注"},   // 作物状态
		"time6":  {Value: task.RemindTime.Format("2006-01-02 15:04")}, // 时间
		"thing7": {Value: "请及时查看您的植物状态"}, // 温馨提示
		"thing5": {Value: task.Content}, // 处理操作
	}

	reqBody := SubscribeMessageRequest{
		ToUser:           task.OpenID,
		TemplateID:       templateID,
		Page:             "/pages/index/index",
		Data:             data,
		MiniprogramState: "formal", // formal, developer, trial
		Lang:             "zh_CN",
	}

	jsonData, _ := json.Marshal(reqBody)
	
	// 在微信云托管中，可以直接调用微信 API 接口，无需 access_token
	apiURL := "http://api.weixin.qq.com/cgi-bin/message/subscribe/send"
	
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("微信接口返回错误: %d - %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// StartRemindWorker 启动定时任务检查器
func StartRemindWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			checkAndSendReminders()
		}
	}()
}

func checkAndSendReminders() {
	ctx := context.Background()
	db := dao.GetDB(ctx)
	if db == nil {
		return
	}
	var tasks []model.RemindTask
	
	// 查找待发送的预约
	now := time.Now()
	err := db.Where("status = ? AND remind_time <= ?", 0, now).Find(&tasks).Error
	if err != nil {
		return
	}

	for _, task := range tasks {
		fmt.Printf("正在发送提醒给 %s, 内容: %s\n", task.OpenID, task.Content)
		err := SendSubscribeMessage(ctx, &task)
		if err != nil {
			fmt.Printf("发送提醒失败 [%d]: %v\n", task.ID, err)
			db.Model(&task).Update("status", 2) // 标记为失败
			continue
		}
		db.Model(&task).Update("status", 1) // 标记为成功
	}
}
