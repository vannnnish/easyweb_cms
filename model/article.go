/**
 * Created by angelina on 2017/9/7.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"gitlab.yeeyuntech.com/yee/easyweb"
	"errors"
	"time"
	"github.com/yeeyuntech/yeego/yeeTime"
)

// 文章
type Article struct {
	Id          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`      // 主键
	CateId      int    `gorm:"not null" json:"cate_id"`                   // 栏目id
	Title       string `gorm:"size:100;not null" json:"title"`            // 标题
	Thumb       string `gorm:"size:100;not null" json:"thumb"`            // 缩略图
	Source      string `gorm:"size:50;not null" json:"source"`            // 来源
	Author      string `gorm:"size:50;not null" json:"author"`            // 作者
	PicAuthor   string `gorm:"size:50;not null" json:"pic_author"`        // 图片作者
	Description string `gorm:"not null" json:"description"`               // 描述
	Content     string `gorm:"type:text;not null" json:"content"`         // 内容
	IsPublish   bool   `gorm:"not null" json:"is_publish" `               // 是否发布
	CreateTime  string `gorm:"type:datetime;not null" json:"create_time"` // 创建时间
	UpdateTime  string `gorm:"type:datetime;not null" json:"update_time"` // 更新或者发布时间
	Sort        int    `gorm:"not null" json:"sort"`                      // 排序
	Hits        int    `gorm:"not null" json:"hits"`                      // 点击量
}

func (Article) TableName() string {
	return "yeecms_article"
}

// 查找文章
func (Article) SelectOne(article *Article) error {
	err := defaultDB.First(article, article.Id).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		err = SelectError
	}
	return nil
}

// 查找文章数量
func (Article) SelectAllCount(cateId int, keyword, source, author string,
	isPublish int, startTime, endTime string) int {
	var count int
	if startTime == "" {
		startTime = "0000-00-00 00:00:00"
	}
	if endTime == "" {
		endTime = "2100-00-00 00:00:00"
	}
	db := defaultDB.Model(&Article{}).
		Where("cate_id = ?", cateId).
		Where("update_time BETWEEN ? AND ?", startTime, endTime)
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if source != "" {
		db = db.Where("source = ?", source)
	}
	if author != "" {
		db = db.Where("author = ?", author)
	}
	if isPublish != 2 {
		db = db.Where("is_publish = ?", isPublish)
	}
	db.Count(&count)
	return count
}

// 分页获取文章数据
func (Article) SelectAll(cateId int, keyword, source, author string,
	isPublish int, startTime, endTime string, pageSize, offset int) []Article {
	var articles []Article
	if startTime == "" {
		startTime = "0000000000"
	}
	if endTime == "" {
		endTime = "9999999999"
	}
	db := defaultDB.Model(&Article{}).
		Where("cate_id = ?", cateId).
		Where("unix_timestamp(update_time) BETWEEN ? AND ?", startTime, endTime).
		Limit(pageSize).Offset(offset).Order("sort DESC ,update_time DESC")
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if source != "" {
		db = db.Where("source = ?", source)
	}
	if author != "" {
		db = db.Where("author = ?", author)
	}
	if isPublish != 2 {
		db = db.Where("is_publish = ?", isPublish)
	}
	db.Find(&articles)
	return articles
}

// 新建文章
func (Article) Create(article Article) error {
	category := &Category{Id: article.CateId}
	if err := (Category{}).SelectOne(category); err != nil {
		return err
	}
	cateModel := &Model{Id: category.ModelId}
	if err := (Model{}).SelectOneModel(cateModel); err != nil {
		return err
	}
	if cateModel.DbTableName != (Article{}).TableName() {
		return errors.New("栏目对应的模型类型错误")
	}
	article.CreateTime = time.Now().Format(yeeTime.FormatMysql)
	if article.UpdateTime == "" {
		article.UpdateTime = time.Now().Format(yeeTime.FormatMysql)
	}
	if err := defaultDB.Create(&article).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return CreateError
	}
	return nil
}

// 编辑文章
func (Article) Update(article Article) error {
	if article.UpdateTime == "" {
		article.UpdateTime = time.Now().Format(yeeTime.FormatMysql)
	}
	if err := defaultDB.Model(article).Update(&article); err != nil {
		easyweb.Logger.Error(err.Error)
		return UpdateError
	}
	return nil
}

// 排序文章
func (Article) DoSort(id, sort int) error {
	if err := defaultDB.Model(&Article{Id: id}).Updates(Article{Sort: sort}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return SortError
	}
	return nil
}

// 发布文章
func (Article) Publish(id int) error {
	if err := defaultDB.Model(&Article{Id: id}).Updates(Article{IsPublish: true}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("发布文章失败")
	}
	return nil
}

// 取消发布文章
func (Article) UnPublish(id int) error {
	if err := defaultDB.Model(&Article{Id: id}).Updates(Article{IsPublish: false}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return errors.New("取消发布文章失败")
	}
	return nil
}

// 删除文章
func (Article) Delete(id int) error {
	if err := defaultDB.Delete(&Article{Id: id}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return DeleteError
	}
	return nil
}

// 访问量+1
func (Article) Hit(id int) {
	article := &Article{Id: id}
	Article{}.SelectOne(article)
	defaultDB.Model(&Article{Id: id}).Updates(Article{Hits: article.Hits + 1})
}
