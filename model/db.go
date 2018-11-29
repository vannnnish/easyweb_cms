/**
 * Created by angelina on 2017/9/7.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var defaultDB *gorm.DB

var (
	SelectError = errors.New("查找失败")
	CreateError = errors.New("新建失败")
	UpdateError = errors.New("更新失败")
	SortError   = errors.New("排序失败")
	DeleteError = errors.New("删除失败")
)

func InitDB(db *gorm.DB) {
	defaultDB = db
}

type DbConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	DbName   string
}

// 初始化测试数据库
func CreateAndInitTestDB() error {
	conf := DbConfig{
		UserName: "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "gorm",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName, conf.Password, conf.Host, conf.Port)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARSET utf8mb4 COLLATE utf8mb4_bin", conf.DbName)
	err = db.Exec(sql).Error
	if err != nil {
		return err
	}
	db.CreateTable()
	db.Close()
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName, conf.Password, conf.Host, conf.Port, conf.DbName)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.LogMode(true)
	defaultDB = db
	return nil
}
