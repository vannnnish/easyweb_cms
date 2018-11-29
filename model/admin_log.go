/**
 * Created by angelina on 2017/9/8.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"time"
	"github.com/yeeyuntech/yeego/yeeTime"
)

// 管理员登录日志
type AdminLog struct {
	Id      int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"` // 主键
	Account string `gorm:"size:100;not null" json:"account"`     // 账号
	Content string `gorm:"size:20;not null" json:"content"`      // 登录或者登出
	IP      string `gorm:"size:20;not null" json:"ip"`           // 登录ip
	Time    string `gorm:"type:datetime;not null" json:"time"`   // 时间
}

func (AdminLog) TableName() string {
	return "yeecms_admin_log"
}

// 记录日志
func (AdminLog) Log(account, ip, content string) {
	log := AdminLog{Account: account, Content: content, IP: ip, Time: time.Now().Format(yeeTime.FormatMysql)}
	defaultDB.Create(&log)
}

// 分页获取日志
func (AdminLog) SelectAll(pageSize, offset int) []AdminLog {
	var logs []AdminLog
	defaultDB.Limit(pageSize).Offset(offset).Order("time DESC").Find(&logs)
	return logs
}

// 获取全部日志数量
func (AdminLog) SelectAllCount() int {
	var count int
	defaultDB.Model(&AdminLog{}).Count(&count)
	return count
}
