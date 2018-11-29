/**
 * Created by angelina on 2017/9/16.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import "gitlab.yeeyuntech.com/yee/easyweb"

// 轮换图
type Slide struct {
	Id          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`  //
	CateId      int    `gorm:"not null" json:"cate_id"`               // 所属栏目id
	Title       string `gorm:"not null" json:"title"`                 // 标题
	Thumb       string `gorm:"not null" json:"thumb"`                 // 缩略图地址
	Url         string `gorm:"not null" json:"url"`                   // 链接地址
	Description string `gorm:"type:text;not null" json:"description"` // 描述
	Sort        int    `gorm:"not null" json:"sort"`                  // 排序
}

func (Slide) TableName() string {
	return "yeecms_slide"
}

// 查找单条轮换图
func (Slide) SelectOne(slide *Slide) error {
	if err := defaultDB.First(slide, slide.Id).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 获取某栏目下轮换图数量
func (Slide) SelectAllCount(cateId int) int {
	var count int
	defaultDB.Model(&Slide{}).Where("cate_id = ?", cateId).Count(&count)
	return count
}

// 分页获取某栏目下轮换图列表
func (Slide) SelectAll(cateId, pageSize, offset int) []Slide {
	var slides []Slide
	defaultDB.Where("cate_id = ?", cateId).Limit(pageSize).Offset(offset).Order("sort DESC").Find(&slides)
	return slides
}

// 新建轮换图
func (Slide) Create(slide Slide) error {
	if err := defaultDB.Create(&slide).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return CreateError
	}
	return nil
}

// 更新
func (Slide) Update(slide Slide) error {
	if err := defaultDB.Model(slide).Update(&slide).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return UpdateError
	}
	return nil
}

// 排序
func (Slide) DoSort(id, sort int) error {
	slide := &Slide{Id: id}
	if err := defaultDB.Model(slide).Updates(Slide{Sort: sort}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return SortError
	}
	return nil
}

// 删除
func (Slide) Delete(id int) error {
	if err := defaultDB.Delete(&Slide{Id: id}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return DeleteError
	}
	return nil
}
