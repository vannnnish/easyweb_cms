/**
 * Created by angelina on 2017/9/11.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"github.com/vannnnish/easyweb_cms/conf"
	"github.com/vannnnish/yeego/yeestrconv"
)

// 管理员角色和权限的中间表，链接了角色以及角色拥有的权限
type AdminRolePrivilege struct {
	Id          int `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	RoleId      int `gorm:"not null" json:"role_id"`      // 角色id
	PrivilegeId int `gorm:"not null" json:"privilege_id"` // 权限id
}

func (AdminRolePrivilege) TableName() string {
	return "yeecms_admin_role_privilege"
}

// 获取某个角色全部操作权限
func (AdminRolePrivilege) RoleActions(roleId int) map[string]string {
	rows := getRoleActions(roleId)
	data := make(map[string]string)
	for _, v := range rows {
		data[yeestrconv.FormatInt(v.CateId)] = v.RoleActions
	}
	return data
}

// 获取某个角色全部操作权限id
func (AdminRolePrivilege) RoleActionsId(roleId int) map[string]string {
	rows := getRoleActions(roleId)
	data := make(map[string]string)
	for _, v := range rows {
		data[yeestrconv.FormatInt(v.CateId)] = v.RoleActionsId
	}
	return data
}

// 临时struct
type TempCateAction struct {
	CateId        int    `gorm:"column:cate_id"`
	RoleActions   string `gorm:"column:role_actions"`
	RoleActionsId string `gorm:"column:role_actions_id"`
}

func getRoleActions(roleId int) []TempCateAction {
	var data []TempCateAction
	if roleId == conf.SuperAdminRoleId {
		// 超级管理员
		defaultDB.Table(AdminPrivilege{}.TableName()).Select("cate_id," +
			"group_concat(yeecms_privilege.action) AS role_actions," +
			"group_concat(yeecms_privilege.id) AS role_actions_id").Group("cate_id").Scan(&data)
	} else {
		defaultDB.Table(AdminPrivilege{}.TableName()).Select("cate_id," +
			"group_concat(yeecms_privilege.action) AS role_actions," +
			"group_concat(yeecms_privilege.id) AS role_actions_id").
			Joins("inner join yeecms_admin_role_privilege " +
				"on yeecms_admin_role_privilege.privilege_id = yeecms_privilege.id").
			Where("yeecms_admin_role_privilege.role_id = ?", roleId).Group("cate_id").Scan(&data)
	}
	return data
}
