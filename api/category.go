/**
 * Created by angelina on 2017/9/13.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"encoding/json"
	"github.com/yeeyuntech/yeego/yeeStrconv"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/conf"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
	"os"
)

type Category_Api struct {
	CategoryModel       model.Category
	AdminPrivilegeModel model.AdminPrivilege
}

// All
// 获取全部栏目以及权限
func (cate Category_Api) All() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		data := cate.CategoryModel.SelectAll(conf.AdminTopCateId, conf.SuperAdminRoleId, true)
		c.Success(data)
	}
	return f
}

// List
// 传入栏目id，获取子栏目信息
func (cate Category_Api) List() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.Param("id").SetDefaultInt(conf.ContentCateId).GetInt()
		cates := cate.CategoryModel.SelectAll(cateId, -1, true)
		c.RenderData("data", cates)
		c.Success(c.GetRenderData())
	}
	return f
}

// Profile
// 单个栏目信息
func (cate Category_Api) Profile() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		category := model.Category{Id: cateId}
		cate.CategoryModel.SelectOne(&category)
		var m map[string]interface{}
		d, _ := json.Marshal(category)
		json.Unmarshal(d, &m)
		m["actions"] = cate.CategoryModel.GetRoleActions(cateId)
		c.Success(m)
	}
	return f
}

// Create
// 创建栏目
func (cate Category_Api) Create() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		name := c.Param("name").MustGetString()
		parentId := c.Param("parent_id").MustGetInt()
		modelId := c.Param("model_id").MustGetInt()
		privileges := c.Param("privileges").GetString()
		contains := c.Param("contains").GetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		category := model.Category{
			Name:     name,
			ParentId: parentId,
			ModelId:  modelId,
			Contain:  contains,
		}
		cateId, err := cate.CategoryModel.Create(category)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := cate.AdminPrivilegeModel.CreateOrUpdate(cateId, privileges); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		os.RemoveAll(conf.CategoryFilePath)
		c.SuccessWithMsg("创建栏目成功")
	}
	return f
}

// Update
// 编辑栏目
// 只能编辑名称 可执行操作
func (cate Category_Api) Update() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.Param("id").MustGetInt()
		name := c.Param("name").MustGetString()
		privileges := c.Param("privileges").GetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := cate.CategoryModel.Update(cateId, name); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := cate.AdminPrivilegeModel.CreateOrUpdate(cateId, privileges); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		os.RemoveAll(conf.CategoryFilePath)
		c.SuccessWithMsg("编辑栏目成功")
	}
	return f
}

// Sort
// 排序栏目
func (cate Category_Api) Sort() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		postData := c.Request().PostForm
		for k, v := range postData {
			err := cate.CategoryModel.DoSort(yeeStrconv.AtoIDefault0(k), yeeStrconv.AtoIDefault0(v[0]))
			if err != nil {
				c.FailWithDefaultCode(err.Error())
				return
			}
		}
		os.RemoveAll(conf.CategoryFilePath)
		c.SuccessWithMsg("排序成功")
	}
	return f
}

// Delete
// 删除栏目
func (cate Category_Api) Delete() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := cate.CategoryModel.Delete(cateId); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		os.RemoveAll(conf.CategoryFilePath)
		c.SuccessWithMsg("删除栏目成功")
	}
	return f
}
