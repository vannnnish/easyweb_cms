/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"errors"
	"github.com/vannnnish/yeego/yeeTime"
	"github.com/vannnnish/easyweb"
	"time"
)

// 后台管理员token
type AdminUserToken struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`      // 主键
	UserId     int    `gorm:"not null;unique" json:"user_id"`            // 用户id
	Token      string `gorm:"not null;size:100" json:"token"`            // 用户token
	UpdateTime string `gorm:"type:datetime;not null" json:"update_time"` // 上次更新token时间
}

func (AdminUserToken) TableName() string {
	return "yeecms_adminuser_token"
}

// 创建或者更新token
func (AdminUserToken) CreateOrUpdate(userId int, token string) error {
	adminUserToken := &AdminUserToken{}
	ok := defaultDB.Where("user_id = ?", userId).First(adminUserToken).RecordNotFound()
	if ok {
		// 不存在
		adminUserToken.UserId = userId
		adminUserToken.Token = token
		adminUserToken.UpdateTime = time.Now().Format(yeeTime.FormatMysql)
		err := defaultDB.Create(adminUserToken).Error
		if err != nil {
			easyweb.Logger.Error(err.Error())
			return errors.New("创建token失败")
		}
		return nil
	}
	// 存在
	adminUserToken.Token = token
	adminUserToken.UpdateTime = time.Now().Format(yeeTime.FormatMysql)
	err := defaultDB.Save(adminUserToken).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("更新token失败")
	}
	return nil
}

// 通过token查找
func (AdminUserToken) SelectOneByToken(token *AdminUserToken) error {
	err := defaultDB.Where(&AdminUserToken{Token: token.Token}).First(token).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 清除用户token
func (AdminUserToken) Clear(userId int) {
	token := &AdminUserToken{UserId: userId}
	defaultDB.Delete(token)
}
