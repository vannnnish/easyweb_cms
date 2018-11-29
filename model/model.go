/**
 * Created by angelina on 2017/9/13.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"errors"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"github.com/yeeyuntech/yeego/yeeStrconv"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/conf"
	"github.com/yeeyuntech/yeego/yeeCache"
	"time"
	"encoding/json"
)

// 模型信息
type Model struct {
	Id          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`  //
	Name        string `gorm:"size:50;not null" json:"name"`          // 模型名称
	Description string `gorm:"not null" json:"description"`           // 描述
	IsShow      bool   `gorm:"not null" json:"is_show"`               // 是否可以选择
	DbTableName string `gorm:"size:50;not null" json:"db_table_name"` // 对应的数据库名称
	Actions     string `gorm:"size:100;not null" json:"actions"`      // 模型拥有的动作权限
}

func (Model) TableName() string {
	return "yeecms_model"
}

// 新建模型
func (Model) Create(model Model) error {
	if !defaultDB.Where("name = ?", model.Name).First(&Model{}).RecordNotFound() {
		return errors.New("模型名称已经存在")
	}
	if err := defaultDB.Create(&model).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return CreateError
	}
	return nil
}

// 获取一个模型信息
func (Model) SelectOneModel(model *Model) error {
	if err := defaultDB.First(model, model.Id).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("获取模型失败")
	}
	return nil
}

// 获取一个模型信息(文件ttl缓存)
func (Model) SelectOneModelWithCache(model *Model) error {
	cacheFileName := conf.CacheFilePath + yeeStrconv.FormatInt(model.Id) + ".cache"
	b, err := yeeCache.FileTtlCache(cacheFileName, func() (b []byte, ttl time.Duration, err error) {
		err = Model{}.SelectOneModel(model)
		if err != nil {
			return nil, time.Second, err
		}
		b, err = json.Marshal(model)
		if err != nil {
			return nil, time.Second, err
		}
		ttl = 24 * time.Hour
		return
	})
	if err != nil {
		if err := (Model{}).SelectOneModel(model); err != nil {
			return err
		}
		return nil
	}
	if err := json.Unmarshal(b, model); err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("获取模型信息失败")
	}
	return nil
}

// 获取全部的模型信息
func (Model) SelectAll() []Model {
	var models []Model
	defaultDB.Model(&Model{}).Find(&models)
	return models
}

// 获取全部的模型信息(文件ttl缓存)
func (Model) SelectAllWithCache() []Model {
	var models []Model
	cacheFileName := conf.CacheFilePath + "all.cache"
	b, err := yeeCache.FileTtlCache(cacheFileName, func() (b []byte, ttl time.Duration, err error) {
		data := Model{}.SelectAll()
		b, err = json.Marshal(&data)
		if err != nil {
			return nil, time.Second, err
		}
		ttl = time.Hour * 24
		return
	})
	if err != nil {
		return Model{}.SelectAll()
	}
	if err := json.Unmarshal(b, &models); err != nil {
		easyweb.Logger.Error(err.Error())
		return nil
	}
	return models
}
