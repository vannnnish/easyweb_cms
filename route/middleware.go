/**
 * Created by angelina on 2017/10/21.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package route

import (
	"github.com/yeeyuntech/yeego"
	"github.com/yeeyuntech/yeego/yeeCrypto"
	"github.com/yeeyuntech/yeego/yeeStrconv"
	"github.com/yeeyuntech/yeego/yeeTransform"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/api"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/conf"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
	"net/http"
	"strings"
	"time"
)

// LoginStateMiddleware
// 登录状态验证中间件
func LoginStateMiddleware() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		loginCookie, err := c.Cookie("yeecms_" + yeego.Config.GetString("app.CookieName"))
		if err != nil || loginCookie == nil {
			c.Fail(1, "登录验证失败")
			c.Abort()
			return
		}
		cookieArray := strings.Split(loginCookie.Value, "|")
		if len(cookieArray) != 3 {
			c.SetCookie(&http.Cookie{
				Name:  loginCookie.Name,
				Value: "",
				Path:  "/",
			})
			c.Fail(1, "登录验证失败")
			c.Abort()
			return
		}
		token := model.AdminUserToken{Token: loginCookie.Value}
		err = model.AdminUserToken{}.SelectOneByToken(&token)
		if err != nil {
			c.SetCookie(&http.Cookie{
				Name:  loginCookie.Name,
				Value: "",
				Path:  "/",
			})
			c.Fail(1, "登录验证失败")
			c.Abort()
			return
		}
		user := model.AdminUser{Id: yeeStrconv.AtoIDefault0(cookieArray[0])}
		err = model.AdminUser{}.SelectOneById(&user)
		if err != nil {
			c.SetCookie(&http.Cookie{
				Name:  loginCookie.Name,
				Value: "",
				Path:  "/",
			})
			c.Fail(1, "登录验证失败")
			c.Abort()
			return
		}
		if user.Id != token.UserId {
			c.SetCookie(&http.Cookie{
				Name:  loginCookie.Name,
				Value: "",
				Path:  "/",
			})
			c.Fail(1, "登录验证失败")
			c.Abort()
			return
		}
		if yeeCrypto.Md5Hex([]byte(user.Password)) != cookieArray[1] {
			c.SetCookie(&http.Cookie{
				Name:  loginCookie.Name,
				Value: "",
				Path:  "/",
			})
			c.Fail(1, "登录验证失败")
			c.Abort()
			return
		}
		if c.Request().Method == "POST" {
			c.Request().ParseForm()
		}
		c.StoreSet("adminUser", yeeTransform.StructToMap(user))
		c.StoreSet("adminUserId", user.Id)
		c.SetCookie(&http.Cookie{
			Name:    loginCookie.Name,
			Value:   loginCookie.Value,
			Path:    "/",
			Expires: time.Now().Add(4 * time.Hour),
		})
		c.Next()
	}
	return f
}

// 执行动作前判断栏目模型与所属模型是否匹配
func ActionPrepareMiddleware() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.PathParam("cateId")
		modelId := c.PathParam("modelId")
		if cateId == "" || modelId == "" || cateId == "0" || modelId == "0" {
			c.FailWithDefaultCode("参数错误")
			c.Abort()
			return
		}
		category := model.Category{Id: yeeStrconv.AtoIDefault0(cateId)}
		err := model.Category{}.SelectOneWithCache(&category)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			c.Abort()
			return
		}
		if category.ModelId != yeeStrconv.AtoIDefault0(modelId) && c.Param("pid").GetInt() != 0 {
			c.FailWithDefaultCode("栏目类型与模型不匹配")
			c.Abort()
			return
		}
		c.Next()
	}
	return f
}

// 根据传入的动作判断是否具有这个权限
func ActionMiddleware() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.PathParam("cateId")
		action := c.PathParam("action")
		if cateId == "" || action == "" || cateId == "0" {
			c.FailWithDefaultCode("参数错误")
			c.Abort()
			return
		}
		category := model.Category{Id: yeeStrconv.AtoIDefault0(cateId)}
		err := model.Category{}.SelectOneWithCache(&category)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			c.Abort()
			return
		}
		m := model.Model{Id: category.ModelId}
		err = model.Model{}.SelectOneModelWithCache(&m)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			c.Abort()
			return
		}
		if !strings.Contains(m.Actions, action) {
			c.Next()
			return
		}
		roleId := c.StoreGetMapInterface("adminUser")["RoleId"].(int)
		actionMap := model.AdminRolePrivilege{}.RoleActions(roleId)
		if value, ok := actionMap[cateId]; ok {
			if strings.Contains(value, action) {
				c.Next()
				return
			}
		}
		c.FailWithDefaultCode("没有权限进行这个操作")
		c.Abort()
	}
	return f
}

var DefaultDispatchMap map[string]map[int]easyweb.HandlerFunc = map[string]map[int]easyweb.HandlerFunc{
	"list":    api.ListModelApiMap,
	"profile": api.ProfileModelApiMap,
	"create":  api.CreateModelApiMap,
	"update":  api.UpdateModelApiMap,
	"sort":    api.SortModelApiMap,
	"publish": api.PublishModelApiMap,
	"delete":  api.DeleteModelApiMap,
}

// 分发中间件
func DispatchMiddleware(dispatchMap map[string]map[int]easyweb.HandlerFunc) easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.PathParam("cateId")
		action := c.PathParam("action")
		modelId := c.PathParam("modelId")
		if cateId == "" || action == "" || modelId == "" || cateId == "0" || modelId == "0" {
			c.FailWithDefaultCode("参数错误")
			c.Abort()
			return
		}
		method := c.Request().Method
		if action == "list" || action == "profile" {
			if method != "GET" {
				c.FailWithDefaultCode("Method Not Allowed")
				c.Abort()
				return
			}
		} else {
			if method != "POST" {
				c.FailWithDefaultCode("Method Not Allowed")
				c.Abort()
				return
			}
		}
		if action == "list" && modelId == yeeStrconv.FormatInt(conf.DirModelId) {
			// 如果是目录类型，则返回子栏目信息
			subCategories := api.Base{}.Category.CategoryModel.
				SelectAllWithCache(yeeStrconv.AtoIDefault0(cateId), c.StoreGetMapInterface("adminUser")["RoleId"].(int), false)
			c.Success(subCategories)
			c.Abort()
			return
		}
		if len(dispatchMap) == 0 {
			dispatchMap = DefaultDispatchMap
		}
		if action == model.Sort {
			api.CommonApiCtx.Sort()(c)
			c.Abort()
			return
		}
		apiMap, ok := dispatchMap[action]
		if !ok {
			c.FailWithDefaultCode("404 NOT FOUND")
			c.Abort()
			return
		}
		function, ok := apiMap[yeeStrconv.AtoIDefault0(modelId)]
		if !ok {
			c.Next()
		} else {
			function(c)
			c.Abort()
		}
	}
	return f
}
