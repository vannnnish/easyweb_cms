/**
 * Created by angelina on 2017/9/11.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/easyweb_cms/model"
	"github.com/vannnnish/yeego/yeestrconv"
)

type AdminRole_Api struct {
	AdminRoleModel      model.AdminRole
	AdminPrivilegeModel model.AdminPrivilege
}

// 返回全部的栏目的可操作信息
func (adminRole AdminRole_Api) AllActions() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		actions := adminRole.AdminPrivilegeModel.GetAllActions()
		c.RenderData("data", actions)
		c.Success(c.GetRenderData())
	}
	return f
}

// 获取全部的角色信息
func (adminRole AdminRole_Api) All() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		count := adminRole.AdminRoleModel.SelectAllCount()
		data := adminRole.AdminRoleModel.SelectAll(100, 0)
		c.RenderData("data", data)
		c.RenderData("count", count)
		c.Success(c.GetRenderData())
	}
	return f
}

// 角色信息列表
func (adminRole AdminRole_Api) List() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		pageSize := c.Param("pagesize").SetDefaultInt(DefaultPageSize).GetInt()
		page := c.Param("page").SetDefaultInt(DefaultPage).GetInt()
		count := adminRole.AdminRoleModel.SelectAllCount()
		data := adminRole.AdminRoleModel.SelectAllId(pageSize, (page-1)*pageSize)
		c.RenderData("data", data)
		c.RenderData("count", count)
		c.Success(c.GetRenderData())
	}
	return f
}

// 获取自己的角色信息
func (adminRole AdminRole_Api) Self() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		roleId := c.StoreGetMapInterface("adminUser")["RoleId"].(int)
		info := adminRole.AdminRoleModel.SelectSelf(roleId)
		c.Success(info)
	}
	return f
}

// 获取某个角色信息
func (adminRole AdminRole_Api) Profile() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		roleId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		info := adminRole.AdminRoleModel.SelectOne(roleId)
		c.Success(info)
	}
	return f
}

// 新建角色
func (adminRole AdminRole_Api) Create() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		name := c.Param("name").MustGetString()
		privIds := c.Param("priv_ids").MustGetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := adminRole.AdminRoleModel.Create(name, privIds); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("新建角色成功")
	}
	return f
}

// 编辑角色
func (adminRole AdminRole_Api) Update() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		roleId := c.Param("id").MustGetInt()
		name := c.Param("name").MustGetString()
		privIds := c.Param("priv_ids").MustGetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := adminRole.AdminRoleModel.Update(roleId, name, privIds); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("编辑角色成功")
	}
	return f
}

// 排序角色
func (adminRole AdminRole_Api) Sort() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		postData := c.Request().PostForm
		for k, v := range postData {
			err := adminRole.AdminRoleModel.DoSort(yeestrconv.AtoIDefault0(k), yeestrconv.AtoIDefault0(v[0]))
			if err != nil {
				c.FailWithDefaultCode(err.Error())
				return
			}
		}
		c.SuccessWithMsg("排序成功")
	}
	return f
}

// 删除角色
func (adminRole AdminRole_Api) Delete() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		roleId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := adminRole.AdminRoleModel.Delete(roleId); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("删除角色成功")
	}
	return f
}
