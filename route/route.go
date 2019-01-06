/**
 * Created by angelina on 2017/9/8.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package route

import (
	"easyweb_cms/api"
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/yeego"
)

func InitRoute(web *easyweb.EasyWeb, dispatchMap map[string]map[int]easyweb.HandlerFunc) {
	//web.Static("/data", "data")
	//web.Static("/static", "view/static")
	var api api.Base
	// Ueditor
	//web.Any("/ueditor", api.Ueditor.UE())
	// 图片上传
	web.POST("/upload/image",
		api.Upload.UploadImage(yeego.Config.GetString("upload.ImgPath"), yeego.Config.GetString("upload.ImgReturnPath")))
	// 文件上传
	web.POST("/upload/file",
		api.Upload.UploadFile(yeego.Config.GetString("upload.FilePath"), yeego.Config.GetString("upload.FileReturnPath")))
	// 视频上传
	web.POST("/upload/video", api.Upload.UploadVideo())
	// 背景图片地址
	web.GET("/bgimage", api.AdminLogin.BgImage())
	// 后台用户登录
	web.POST("/login", api.AdminLogin.Login())
	// 清除缓存
	web.GET("/clearcache", api.AdminLogin.ClearCache())
	// 后台组
	adminGroup := web.Group("/admin", LoginStateMiddleware())
	// 后台用户退出登录
	adminGroup.GET("/logout", api.AdminLogin.Logout())
	// 登录用户的信息
	adminGroup.GET("/userinfo", api.AdminUser.Self())
	// 全部的栏目以及可执行操作和当前用户可执行操作信息
	adminGroup.GET("/category/all", api.Category.All())
	// 该用户角色的操作
	adminGroup.GET("/role/actions", api.AdminRole.Self())
	adminGroup.GET("/model/all", api.Model.List())
	// 全部的角色信息
	adminGroup.GET("/roles", api.AdminRole.All())
	// 所有的操作信息
	adminGroup.GET("/actions", api.AdminRole.AllActions())
	// 重置管理员密码
	adminGroup.POST("/adminuser/resetpwd", api.AdminUser.ResetPassword())
	// 管理员修改自己的密码
	adminGroup.POST("/adminuser/changepwd", api.AdminUser.UpdatePassword())
	// 一些分发函数
	adminGroup.GET("/api/:action/:cateId/:modelId", api.DefaultPage(),
		ActionPrepareMiddleware(),
		ActionMiddleware(),
		DispatchMiddleware(dispatchMap))
	adminGroup.POST("/api/:action/:cateId/:modelId", api.DefaultPage(),
		ActionPrepareMiddleware(),
		ActionMiddleware(),
		DispatchMiddleware(dispatchMap))
	//adminGroup.Any("/api/:action/:cateId/:modelId", api.DefaultPage(),
	//	ActionPrepareMiddleware(),
	//	ActionMiddleware(),
	//	DispatchMiddleware(dispatchMap))
}
