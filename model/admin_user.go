/**
 * Created by angelina on 2017/9/8.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"errors"
	"github.com/yeeyuntech/yeego/yeeCrypto"
	"time"
	"github.com/yeeyuntech/yeego/yeeTime"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"fmt"
)

// 后台管理员
type AdminUser struct {
	Id            int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`    // 主键
	Account       string `gorm:"size:100;not null" json:"account"`        // 账号唯一不可变
	Password      string `gorm:"not null" json:"password"`                // 密码
	UserName      string `gorm:"size:100;not null" json:"user_name"`      // 用户真实姓名
	RoleId        int    `gorm:"not null" json:"role_id"`                 // 角色id
	Phone         string `gorm:"size:20;not null" json:"phone"`           // 电话号码
	Email         string `gorm:"size:50;not null" json:"email"`           // 邮箱地址
	Sort          int    `gorm:"not null" json:"sort"`                    // 排序
	LoginIP       string `gorm:"size:20;not null" json:"login_ip"`        // 登录ip
	LastLoginTime string `gorm:"size:20;not null" json:"last_login_time"` // 上次登录时间
	ThisLoginTime string `gorm:"size:20;not null" json:"this_login_time"` // 本次登录时间
}

func (AdminUser) TableName() string {
	return "yeecms_adminuser"
}

// 管理员登录
func (AdminUser) Login(account, pwd, ip string) error {
	var adminUser AdminUser
	err := defaultDB.Where(&AdminUser{Account: account}).First(&adminUser).Error
	if err != nil && adminUser.Id == 0 {
		return errors.New("账号或者密码错误，请重试")
	}
	if adminUser.Password != yeeCrypto.Sha256Hex([]byte(pwd)) {
		return errors.New("账号或者密码错误，请重试")
	}
	adminUser.LoginIP = ip
	adminUser.LastLoginTime = adminUser.ThisLoginTime
	adminUser.ThisLoginTime = time.Now().Format(yeeTime.FormatMysql)
	defaultDB.Model(&AdminUser{}).Updates(adminUser)
	return nil
}

// 新建管理员用户
func (AdminUser) Create(user AdminUser) error {
	var u AdminUser
	defaultDB.Where(&AdminUser{Account: user.Account}).First(&u)
	if u.Id != 0 {
		return errors.New("账号已经存在，请重试")
	}
	user.Password = yeeCrypto.Sha256Hex([]byte(user.Password))
	err := defaultDB.Create(&user).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return CreateError
	}
	return nil
}

// 通过id查找管理员
func (AdminUser) SelectOneById(user *AdminUser) error {
	err := defaultDB.First(user, user.Id).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 通过account查找管理员
func (AdminUser) SelectOneByAccount(user *AdminUser) error {
	err := defaultDB.Where(&AdminUser{Account: user.Account}).First(user).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 分页获取管理员数据，除去默认账号
func (AdminUser) SelectAllWithoutDefault(condition map[string]string, pageSize, offset int) []AdminUser {
	var users []AdminUser
	db := defaultDB.Model(&AdminUser{}).Limit(pageSize).Offset(offset)
	for k, v := range condition {
		if k != "" {
			if k == "role_id" {
				if v != "0" {
					db = db.Where("role_id = ?", v)
				}
			} else {
				db = db.Where(fmt.Sprintf("%s LIKE ?", k), "%"+v+"%")
			}
		}
	}
	err := db.Where("account <> ?", "yeeyun_root").Order("sort DESC,id DESC").Find(&users).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return nil
	}
	for i := 0; i < len(users); i++ {
		users[i].Password = ""
	}
	return users
}

// 获取管理员数量，除去默认账号
func (AdminUser) SelectAllWithoutDefaultCount(condition map[string]string) int {
	var count int
	db := defaultDB.Model(&AdminUser{}).Where("account <> ?", "yeeyun_root")
	for k, v := range condition {
		if k != "" {
			if k == "role_id" {
				if v != "0" {
					db = db.Where("role_id = ?", v)
				}
			} else {
				db = db.Where(fmt.Sprintf("%s LIKE ?", k), "%"+v+"%")
			}
		}
	}
	err := db.Count(&count).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return 0
	}
	return count
}

// 更新管理员用户信息
func (AdminUser) Update(user AdminUser) error {
	err := defaultDB.Save(&user).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return UpdateError
	}
	return nil
}

// 更新管理员密码
func (AdminUser) UpdatePassword(id int, oldPwd, newPwd string) error {
	user := &AdminUser{Id: id}
	err := AdminUser{}.SelectOneById(user)
	if err != nil {
		return err
	}
	if user.Password != yeeCrypto.Sha256Hex([]byte(oldPwd)) {
		return errors.New("旧密码输入错误，请重试")
	}
	err = defaultDB.Model(user).Updates(AdminUser{Password: yeeCrypto.Sha256Hex([]byte(newPwd))}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("修改密码失败,请重试")
	}
	return nil
}

// 重置管理员密码
func (AdminUser) ResetPassword(id int, pwd string) error {
	user := &AdminUser{Id: id}
	err := AdminUser{}.SelectOneById(user)
	if err != nil {
		return err
	}
	err = defaultDB.Model(user).Updates(AdminUser{Password: yeeCrypto.Sha256Hex([]byte(pwd))}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("重置密码失败,请重试")
	}
	return nil
}

// 排序
func (AdminUser) DoSort(id, sort int) error {
	user := &AdminUser{Id: id}
	err := AdminUser{}.SelectOneById(user)
	if err != nil {
		return err
	}
	err = defaultDB.Model(user).Updates(AdminUser{Sort: sort}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("排序失败,请重试")
	}
	return nil
}

// 删除
func (AdminUser) Delete(id int) error {
	user := &AdminUser{Id: id}
	err := defaultDB.Delete(user).Error
	if err != nil {
		easyweb.Logger.Error(err.Error)
		return DeleteError
	}
	return nil
}
