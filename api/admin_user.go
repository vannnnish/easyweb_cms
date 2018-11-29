/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"github.com/yeeyuntech/yeego"
	"github.com/yeeyuntech/yeego/yeeCrypto"
	"github.com/yeeyuntech/yeego/yeeStrconv"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/conf"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
)

type AdminUser_Api struct {
	AdminUserModel model.AdminUser
}

// 获取管理员个人信息
func (adminUser AdminUser_Api) Self() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		userId := c.StoreGetInt("adminUserId")
		user := &model.AdminUser{Id: userId}
		adminUser.AdminUserModel.SelectOneById(user)
		user.Password = ""
		c.Success(*user)
	}
	return f
}

// 新建管理员
func (adminUser AdminUser_Api) Create() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		account := c.Param("account").MustGetString()
		roleId := c.Param("role_id").MustGetInt()
		defaultPwd := yeeCrypto.Sha256Hex([]byte(yeego.Config.GetString("app.DefaultPassword")))
		password := c.Param("password").SetDefault(defaultPwd).GetString()
		username := c.Param("user_name").GetString()
		phone := c.Param("phone").GetString()
		email := c.Param("email").GetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		user := model.AdminUser{
			Account:  account,
			RoleId:   roleId,
			Password: password,
			UserName: username,
			Phone:    phone,
			Email:    email,
		}
		err := adminUser.AdminUserModel.Create(user)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("新建管理员成功")
	}
	return f
}

// 管理员列表
func (adminUser AdminUser_Api) List() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		pageSize := c.Param("pagesize").SetDefaultInt(DefaultPageSize).GetInt()
		page := c.Param("page").SetDefaultInt(DefaultPage).GetInt()
		roleId := c.Param("role_id").GetInt()
		account := c.Param("account").GetString()
		count := adminUser.AdminUserModel.SelectAllWithoutDefaultCount(map[string]string{
			"role_id": yeeStrconv.FormatInt(roleId),
			"account": account,
		})
		data := adminUser.AdminUserModel.SelectAllWithoutDefault(map[string]string{
			"role_id": yeeStrconv.FormatInt(roleId),
			"account": account,
		}, pageSize, (page-1)*pageSize)
		c.RenderData("data", data)
		c.RenderData("count", count)
		c.Success(c.GetRenderData())
	}
	return f
}

// 获取管理员信息
func (adminUser AdminUser_Api) Profile() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		userId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		user := &model.AdminUser{Id: userId}
		adminUser.AdminUserModel.SelectOneById(user)
		user.Password = ""
		c.Success(*user)
	}
	return f
}

// 更新管理员信息
func (adminUser AdminUser_Api) Update() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		adminUserId := c.Param("id").MustGetInt()
		username := c.Param("user_name").GetString()
		phone := c.Param("phone").GetString()
		email := c.Param("email").GetString()
		roleId := c.Param("role_id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		user := &model.AdminUser{Id: adminUserId}
		err := adminUser.AdminUserModel.SelectOneById(user)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		user.UserName = username
		user.Phone = phone
		user.Email = email
		user.RoleId = roleId
		err = adminUser.AdminUserModel.Update(*user)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("更新管理员用户成功")
	}
	return f
}

// 修改自己的密码
func (adminUser AdminUser_Api) UpdatePassword() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		adminUserId := c.StoreGetInt("adminUserId")
		oldPwd := c.Param("old_pwd").MustGetString()
		newPwd := c.Param("new_pwd").MustGetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := adminUser.AdminUserModel.UpdatePassword(adminUserId, oldPwd, newPwd); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("修改密码成功")
	}
	return f
}

// 重置密码
func (adminUser AdminUser_Api) ResetPassword() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		if c.StoreGetMapInterface("adminUser")["RoleId"].(int) != conf.SuperAdminRoleId {
			c.FailWithDefaultCode("没有权限重置密码")
			return
		}
		adminUserId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		pwd := yeeCrypto.Sha256Hex([]byte(yeego.Config.GetString("app.DefaultPassword")))
		if err := adminUser.AdminUserModel.ResetPassword(adminUserId, pwd); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("重置密码成功，新密码为:" + yeego.Config.GetString("app.DefaultPassword"))
	}
	return f
}

// 排序
func (adminUser AdminUser_Api) Sort() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		postData := c.Request().PostForm
		for k, v := range postData {
			if err := adminUser.AdminUserModel.DoSort(yeeStrconv.AtoIDefault0(k), yeeStrconv.AtoIDefault0(v[0])); err != nil {
				c.FailWithDefaultCode(err.Error())
				return
			}
		}
		c.SuccessWithMsg("排序成功")
	}
	return f
}

// 删除
func (adminUser AdminUser_Api) Delete() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		delAdminUserId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := adminUser.AdminUserModel.Delete(delAdminUserId); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("删除用户成功")
	}
	return f
}
