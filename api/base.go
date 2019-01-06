/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"easyweb_cms/conf"
	"github.com/vannnnish/easyweb"
)

const (
	// 默认的分页每页数量
	DefaultPageSize int = 20
	DefaultPage     int = 1
)

type Base struct {
	Common     CommonApi
	Ueditor    Ueditor_Api
	Upload     Upload_Api
	Category   Category_Api
	Model      Model_Api
	AdminLogin AdminLogin_Api
	AdminLog   AdminLog_Api
	AdminRole  AdminRole_Api
	AdminUser  AdminUser_Api
}

func (Base) DefaultPage() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		c.FailWithDefaultCode("404 NOT FOUND")
	}
	return f
}

var ListModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.AdminRoleModelId:    AdminRole_Api{}.List(),
	conf.AdminUserLogModelId: AdminLog_Api{}.List(),
	conf.AdminUserModelId:    AdminUser_Api{}.List(),
	conf.CategoryModelId:     Category_Api{}.List(),
	conf.ArticleModelId:      Article_Api{}.List(),
	//conf.SlideModelId:        Slide_Api{}.List(),
	//conf.DownloadModelId:     Download_Api{}.List(),
}

var ProfileModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.AdminRoleModelId: AdminRole_Api{}.Profile(),
	conf.AdminUserModelId: AdminUser_Api{}.Profile(),
	conf.CategoryModelId:  Category_Api{}.Profile(),
	conf.ArticleModelId:   Article_Api{}.Profile(),
	//conf.SinglePageModelId: SinglePage_Api{}.Profile(),
	//conf.SlideModelId:      Slide_Api{}.Profile(),
	//conf.DownloadModelId:   Download_Api{}.Profile(),
}

var CreateModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.AdminRoleModelId: AdminRole_Api{}.Create(),
	conf.AdminUserModelId: AdminUser_Api{}.Create(),
	conf.CategoryModelId:  Category_Api{}.Create(),
	conf.ArticleModelId:   Article_Api{}.Create(),
	//conf.SlideModelId:     Slide_Api{}.Create(),
	//conf.DownloadModelId:  Download_Api{}.Create(),
}

var UpdateModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.AdminRoleModelId: AdminRole_Api{}.Update(),
	conf.AdminUserModelId: AdminUser_Api{}.Update(),
	conf.CategoryModelId:  Category_Api{}.Update(),
	conf.ArticleModelId:   Article_Api{}.Update(),
	//conf.SinglePageModelId: SinglePage_Api{}.Update(),
	//conf.SlideModelId:      Slide_Api{}.Update(),
	//conf.DownloadModelId:   Download_Api{}.Update(),
}

var SortModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.AdminRoleModelId: AdminRole_Api{}.Sort(),
	conf.AdminUserModelId: AdminUser_Api{}.Sort(),
	conf.CategoryModelId:  Category_Api{}.Sort(),
	conf.ArticleModelId:   Article_Api{}.Sort(),
	//conf.SlideModelId:     Slide_Api{}.Sort(),
	//conf.DownloadModelId:  Download_Api{}.Sort(),
}

var PublishModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.ArticleModelId: Article_Api{}.Publish(),
}

var DeleteModelApiMap map[int]easyweb.HandlerFunc = map[int]easyweb.HandlerFunc{
	conf.AdminRoleModelId: AdminRole_Api{}.Delete(),
	conf.AdminUserModelId: AdminUser_Api{}.Delete(),
	conf.CategoryModelId:  Category_Api{}.Delete(),
	conf.ArticleModelId:   Article_Api{}.Delete(),
	//conf.SlideModelId:     Slide_Api{}.Delete(),
	//conf.DownloadModelId:  Download_Api{}.Delete(),
}
