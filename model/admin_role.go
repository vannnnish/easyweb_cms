/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"gitlab.yeeyuntech.com/yee/easyweb_cms/conf"
	"github.com/yeeyuntech/yeego/yeeSql"
	"github.com/yeeyuntech/yeego/yeeTransform"
	"github.com/yeeyuntech/yeego/yeeStrconv"
	"errors"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"github.com/yeeyuntech/yeego/yeeStrings"
)

// 管理员角色
type AdminRole struct {
	Id   int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"` // 主键
	Name string `gorm:"size:50;not null" json:"name"`         // 角色名称
	Sort int    `gorm:"not null" json:"sort"`                 // 角色排序
}

func (AdminRole) TableName() string {
	return "yeecms_admin_role"
}

// 分页获取全部的角色信息,其中角色的栏目权限用操作表示
func (AdminRole) SelectAll(pageSize, offset int) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	temp["id"] = "-1"
	temp["name"] = "超级管理员"
	temp["sort"] = "0"
	temp["role_cates"] = AdminRolePrivilege{}.RoleActions(conf.SuperAdminRoleId)
	data = append(data, temp)
	if offset > 0 {
		offset = offset - 1
	}
	// 获取全部role信息
	rows, _ := defaultDB.Table(AdminRole{}.TableName()).Limit(pageSize).Offset(offset).Rows()
	otherData := yeeSql.RowsToMapSlice(rows)
	infoOtherInterface := yeeTransform.MapSliceStringToInterface(otherData)
	for k, v := range otherData {
		infoOtherInterface[k]["role_cates"] = AdminRolePrivilege{}.RoleActions(yeeStrconv.AtoIDefault0(v["id"]))
	}
	data = append(data, infoOtherInterface...)
	return data
}

// 分页获取全部的角色信息,其中角色的栏目权限用操作id表示
func (AdminRole) SelectAllId(pageSize, offset int) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	temp["id"] = "-1"
	temp["name"] = "超级管理员"
	temp["sort"] = "0"
	temp["role_cates"] = AdminRolePrivilege{}.RoleActionsId(conf.SuperAdminRoleId)
	data = append(data, temp)
	if offset > 0 {
		offset = offset - 1
	}
	// 获取全部role信息
	rows, _ := defaultDB.Table(AdminRole{}.TableName()).Limit(pageSize).Offset(offset).Rows()
	otherData := yeeSql.RowsToMapSlice(rows)
	infoOtherInterface := yeeTransform.MapSliceStringToInterface(otherData)
	for k, v := range otherData {
		infoOtherInterface[k]["role_cates"] = AdminRolePrivilege{}.RoleActionsId(yeeStrconv.AtoIDefault0(v["id"]))
	}
	data = append(data, infoOtherInterface...)
	return data
}

// 获取全部角色的数量
func (AdminRole) SelectAllCount() int {
	count := 0
	defaultDB.Table(AdminRole{}.TableName()).Count(&count)
	return count + 1
}

// 获取单个角色信息,其中角色的栏目权限用操作id表示
func (AdminRole) SelectOne(roleId int) map[string]interface{} {
	if roleId == conf.SuperAdminRoleId {
		data := make(map[string]interface{})
		data["id"] = "-1"
		data["name"] = "超级管理员"
		data["sort"] = "0"
		data["role_cates"] = AdminRolePrivilege{}.RoleActionsId(roleId)
		return data
	}
	rows, err := defaultDB.Table(AdminRole{}.TableName()).Where("id = ?", roleId).Rows()
	if err != nil {
		return nil
	}
	infos := yeeSql.RowsToMapSlice(rows)
	if len(infos) <= 0 {
		return nil
	}
	info := infos[0]
	infoMap := yeeTransform.MapStringToInterface(info)
	infoMap["role_cates"] = AdminRolePrivilege{}.RoleActionsId(roleId)
	return infoMap
}

// 获取单个角色信息,其中角色的栏目权限用操作表示
func (AdminRole) SelectSelf(roleId int) map[string]interface{} {
	if roleId == conf.SuperAdminRoleId {
		data := make(map[string]interface{})
		data["id"] = "-1"
		data["name"] = "超级管理员"
		data["sort"] = "0"
		data["role_cates"] = AdminRolePrivilege{}.RoleActions(roleId)
		return data
	}
	rows, err := defaultDB.Table(AdminRole{}.TableName()).Where("id = ?", roleId).Rows()
	if err != nil {
		return nil
	}
	infos := yeeSql.RowsToMapSlice(rows)
	if len(infos) <= 0 {
		return nil
	}
	info := infos[0]
	infoMap := yeeTransform.MapStringToInterface(info)
	infoMap["role_cates"] = AdminRolePrivilege{}.RoleActions(roleId)
	return infoMap
}

// 创建角色
func (AdminRole) Create(name, privilegeIds string) error {
	if !defaultDB.Where(AdminRole{Name: name}).First(&AdminRole{}).RecordNotFound() {
		return errors.New("角色名称已经存在，请重试")
	}
	tx := defaultDB.Begin()
	// 新建角色
	newRole := &AdminRole{Name: name}
	err := tx.Create(newRole).Error
	if err != nil {
		tx.Rollback()
		easyweb.Logger.Error(err.Error())
		return CreateError
	}
	// 为角色添加权限
	for _, v := range yeeStrings.StringToIntArray(privilegeIds, ",") {
		err := tx.Create(&AdminRolePrivilege{RoleId: newRole.Id, PrivilegeId: v}).Error
		if err != nil {
			tx.Rollback()
			easyweb.Logger.Error(err.Error())
			return CreateError
		}
	}
	tx.Commit()
	return nil
}

// 编辑角色信息
func (AdminRole) Update(roleId int, name, privilegeIds string) error {
	if roleId == conf.SuperAdminRoleId {
		return errors.New("不允许修改超级管理员")
	}
	if !defaultDB.Where("id <> ?", roleId).Where("name = ?", name).First(&AdminRole{}).RecordNotFound() {
		return errors.New("角色名称已经存在，请重试")
	}
	tx := defaultDB.Begin()
	// 更新name
	err := tx.Model(&AdminRole{Id: roleId}).Updates(AdminRole{Name: name}).Error
	if err != nil {
		tx.Rollback()
		easyweb.Logger.Error(err.Error())
		return errors.New("更新角色信息失败，请重试")
	}
	// 删除原有权限
	err = tx.Where("role_id = ?", roleId).Delete(&AdminRolePrivilege{}).Error
	if err != nil {
		tx.Rollback()
		easyweb.Logger.Error(err.Error())
		return errors.New("更新角色信息失败，请重试")
	}
	// 新建权限
	for _, v := range yeeStrings.StringToIntArray(privilegeIds, ",") {
		err := tx.Create(&AdminRolePrivilege{RoleId: roleId, PrivilegeId: v}).Error
		if err != nil {
			tx.Rollback()
			easyweb.Logger.Error(err.Error())
			return errors.New("新建角色失败，请重试")
		}
	}
	tx.Commit()
	return nil
}

// 排序角色
func (AdminRole) DoSort(id int, sort int) error {
	err := defaultDB.Model(&AdminRole{Id: id}).Updates(AdminRole{Sort: sort}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("排序角色失败")
	}
	return nil
}

// 删除角色
func (AdminRole) Delete(roleId int) error {
	if roleId == conf.SuperAdminRoleId {
		return errors.New("不能删除超级管理员")
	}
	// 判断是否有该角色的管理员
	count := AdminUser{}.SelectAllWithoutDefaultCount(map[string]string{
		"role_id": yeeStrconv.FormatInt(roleId),
	})
	if count > 0 {
		return errors.New("删除角色失败，请先删除此角色下的管理员")
	}
	// 删除角色
	err := defaultDB.Delete(&AdminRole{Id: roleId}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return DeleteError
	}
	// 删除角色和权限的中间表数据
	defaultDB.Where("role_id = ?", roleId).Delete(AdminRolePrivilege{})
	return nil
}
