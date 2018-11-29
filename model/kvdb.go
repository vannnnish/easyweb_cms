/**
 * Created by angelina on 2017/9/16.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import "gitlab.yeeyuntech.com/yee/easyweb"

type Kvdb struct {
	Key   string `gorm:"not null" json:"key"`                 // 键
	Value string `gorm:"type:longblob;not null" json:"value"` // 值
}

func (Kvdb) TableName() string {
	return "yeecms_kvdb"
}

func (Kvdb) Set(kv Kvdb) error {
	if defaultDB.Where("`key` = ?", kv.Key).First(&Kvdb{}).RecordNotFound() {
		// 不存在，创建
		if err := defaultDB.Create(kv).Error; err != nil {
			easyweb.Logger.Error(err.Error())
			return CreateError
		}
		return nil
	}
	// 存在，更新
	if err := defaultDB.Table(Kvdb{}.TableName()).Where("`key` = ?", kv.Key).Updates(Kvdb{Value: kv.Value}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return UpdateError
	}
	return nil
}

func (Kvdb) Get(kv *Kvdb) {
	defaultDB.Where("`key` = ?", kv.Key).First(kv)
}
