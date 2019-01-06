/**
 * Created by angelina on 2017/9/16.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"github.com/vannnnish/yeego/yeeTime"
	"github.com/vannnnish/easyweb"
	"time"
)

// 单页面
type SinglePage struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`      //
	CateId     int    `gorm:"not null" json:"cate_id"`                   // 所属栏目id
	Title      string `gorm:"size:50;not null" json:"title"`             // 标题
	Source     string `gorm:"size:30;not null" json:"source"`            // 来源
	Author     string `gorm:"size:30;not null" json:"author"`            // 作者
	Content    string `gorm:"type:longtext;not null" json:"content"`     // 内容
	UpdateTime string `gorm:"type:datetime;not null" json:"update_time"` // 更新时间
	Hits       int    `gorm:"not null" json:"hits"`                      // 点击量
}

func (SinglePage) TableName() string {
	return "yeecms_singlepage"
}

// 新建或者编辑单页面
func (SinglePage) CreateOrUpdate(singlePage SinglePage) error {
	singlePage.UpdateTime = time.Now().Format(yeeTime.FormatMysql)
	if defaultDB.Where(&SinglePage{CateId: singlePage.CateId}).First(&SinglePage{}).RecordNotFound() {
		// 未找到，创建
		if err := defaultDB.Create(&singlePage).Error; err != nil {
			easyweb.Logger.Error(err.Error())
			return CreateError
		}
		return nil
	}
	// 找到了，更新
	if err := defaultDB.Model(&SinglePage{Id: singlePage.Id}).Updates(singlePage).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return UpdateError
	}
	return nil
}

// 查找
func (SinglePage) SelectOne(singlePage *SinglePage) error {
	if err := defaultDB.Where("cate_id = ?", singlePage.CateId).First(singlePage).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 删除单页面
func (SinglePage) Delete(cateId int) error {
	if err := defaultDB.Where("cate_id = ?", cateId).Delete(&SinglePage{}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return DeleteError
	}
	return nil
}

// 访问量+1
func (SinglePage) Hit(cateId int) {
	singlePage := &SinglePage{CateId: cateId}
	SinglePage{}.SelectOne(singlePage)
	defaultDB.Model(&SinglePage{CateId: cateId}).Updates(SinglePage{Hits: singlePage.Hits + 1})
}
