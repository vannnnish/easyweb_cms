/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"errors"
	"github.com/yeeyuntech/yeego"
	"github.com/yeeyuntech/yeego/yeeStrings"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"strings"
)

// 一些默认的操作
const (
	Create  string = "create"
	Delete  string = "delete"
	Profile string = "profile"
	Sort    string = "sort"
	Update  string = "update"
	Publish string = "publish"
)

// 后台权限
type AdminPrivilege struct {
	Id     int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"` //
	CateId int    `gorm:"not null" json:"cate_id"`              // 栏目id
	Action string `gorm:"size:50;not null" json:"action"`       // 对应的操作
}

func (AdminPrivilege) TableName() string {
	return "yeecms_privilege"
}

// 获取全部的可执行操作string
func getAllActionStr() []string {
	actions := AdminPrivilege{}.GetAllActions()
	data := make([]string, 0)
	for k := range actions {
		data = append(data, k)
	}
	return data
}

// 获取全部操作 {'add':'添加'，‘delete’:' 删除',…}
func (AdminPrivilege) GetAllActions() map[string]string {
	other := yeego.Config.GetStringMapString("actions")
	data := make(map[string]string)
	data[Create] = "添加"
	data[Delete] = "删除"
	data[Profile] = "查看"
	data[Sort] = "排序"
	data[Update] = "编辑"
	data[Publish] = "发布"
	for k, v := range other {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	return data
}

// 是否是正确的action
func (AdminPrivilege) isRightAction(list []string) bool {
	if len(list) == 0 {
		return false
	}
	right := true
	for _, v := range list {
		if !yeeStrings.IsInSlice(getAllActionStr(), v) {
			right = false
			break
		}
	}
	return right
}

// 批量创建/更新某个栏目的权限
// privilege = "add,delete,profile..."
func (AdminPrivilege) CreateOrUpdate(cateId int, privilege string) error {
	if privilege == "" {
		privilege = Profile
	} else {
		privilege += "," + Profile
	}
	tx := defaultDB.Begin()
	// 首先删除掉这个栏目的除了profile以外的其他action
	tx.Where("cate_id = ?", cateId).Where("action <> ?", Profile).Delete(AdminPrivilege{})
	// 判断传入的action是否正确
	actionList := strings.Split(privilege, ",")
	if !(AdminPrivilege{}.isRightAction(actionList)) {
		tx.Rollback()
		return errors.New("传入权限类型错误，请检查")
	}
	// 依次创建栏目权限
	for _, v := range actionList {
		err := tx.Where(AdminPrivilege{CateId: cateId, Action: v}).FirstOrCreate(&AdminPrivilege{}).Error
		if err != nil {
			tx.Rollback()
			easyweb.Logger.Error(err.Error())
			return errors.New("创建栏目权限失败，请检查")
		}
	}
	tx.Commit()
	return nil
}

// 批量删除某个栏目的权限
func (AdminPrivilege) Delete(cateId int) {
	defaultDB.Where(AdminPrivilege{CateId: cateId}).Delete(AdminPrivilege{})
}

// 栏目是否包含这个操作权限
func (AdminPrivilege) IsCateActionExist(cateId int, action string) bool {
	return !defaultDB.Where(AdminPrivilege{CateId: cateId, Action: action}).
		FirstOrCreate(&AdminPrivilege{}).RecordNotFound()
}
