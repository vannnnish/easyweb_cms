/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"github.com/buger/jsonparser"
	"github.com/yeeyuntech/yeego"
	"github.com/yeeyuntech/yeego/yeeCache"
	"github.com/yeeyuntech/yeego/yeeCrypto"
	"github.com/yeeyuntech/yeego/yeeHttp"
	"github.com/yeeyuntech/yeego/yeeStrconv"
	"github.com/yeeyuntech/yeego/yeeTime"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/conf"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
	"net/http"
	"os"
	"strings"
	"time"
)

type AdminLogin_Api struct {
	AdminUserModel          model.AdminUser
	AdminUserTokenModel     model.AdminUserToken
	AdminLogModel           model.AdminLog
	AdminRolePrivilegeModel model.AdminRolePrivilege
}

func cookieName() string {
	return "yeecms_" + yeego.Config.GetString("app.CookieName")
}

// 获取背景图
func (adminLogin AdminLogin_Api) BgImage() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		b, err := yeeCache.FileTtlCache(conf.BgImageCachePath, func() (b []byte, ttl time.Duration, err error) {
			data, err := yeeHttp.Get("http://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1").Exec().ToBytes()
			if err != nil {
				return nil, 24 * time.Hour, err
			}
			url, _, _, err := jsonparser.Get(data, "images", "[0]", "url")
			if err != nil {
				return nil, 24 * time.Hour, err
			}
			return url, 24 * time.Hour, nil
		})
		if err != nil {
			c.Success("/static/bgImage.jpg")
			return
		}
		c.Success("http://cn.bing.com" + string(b))
	}
	return f
}

// 后台管理员登录
func (adminLogin AdminLogin_Api) Login() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		account := c.Param("account").MustGetString()
		password := c.Param("password").MustGetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := adminLogin.AdminUserModel.Login(account, password, c.ClientIP()); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		adminUser := &model.AdminUser{Account: account}
		adminLogin.AdminUserModel.SelectOneByAccount(adminUser)
		token := strings.Join([]string{yeeStrconv.FormatInt(adminUser.Id), yeeCrypto.Md5Hex([]byte(adminUser.Password)),
			yeeTime.TimeToUnixS(time.Now())}, "|")
		if err := adminLogin.AdminUserTokenModel.CreateOrUpdate(adminUser.Id, token); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		cookie := &http.Cookie{
			Name:    cookieName(),
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(4 * time.Hour),
		}
		c.SetCookie(cookie)
		if account != "yeeyun_root" {
			adminLogin.AdminLogModel.Log(account, c.ClientIP(), "登录成功")
		}
		// 获取该角色的权限信息
		privileges := adminLogin.AdminRolePrivilegeModel.RoleActions(adminUser.RoleId)
		c.RenderData("privileges", privileges)
		adminUser.Password = ""
		c.RenderData("user_info", *adminUser)
		c.Success(c.GetRenderData())
	}
	return f
}

// 退出登录
func (adminLogin AdminLogin_Api) Logout() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		adminLogin.AdminUserTokenModel.Clear(c.StoreGetInt("adminUserId"))
		if c.StoreGetMapInterface("adminUser")["Account"].(string) != "yeeyun_root" {
			adminLogin.AdminLogModel.Log(c.StoreGetMapInterface("adminUser")["Account"].(string), c.ClientIP(), "退出登录")
		}
		c.SuccessWithMsg("退出登录成功")
	}
	return f
}

// 清除cache
func (adminLogin AdminLogin_Api) ClearCache() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		os.RemoveAll(conf.CacheFilePath)
		os.RemoveAll(conf.CategoryFilePath)
		c.SuccessWithMsg("清除缓存成功")
	}
	return f
}
