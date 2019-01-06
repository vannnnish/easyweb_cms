/**
 * Created by angelina on 2017/9/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"easyweb_cms/model"
	"github.com/vannnnish/easyweb"
)

type AdminLog_Api struct {
	AdminLogModel model.AdminLog
}

// 日志列表数据
func (adminLog AdminLog_Api) List() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		pageSize := c.Param("pagesize").SetDefaultInt(DefaultPageSize).GetInt()
		page := c.Param("page").SetDefaultInt(DefaultPage).GetInt()
		info := adminLog.AdminLogModel.SelectAll(pageSize, (page-1)*pageSize)
		count := adminLog.AdminLogModel.SelectAllCount()
		c.RenderData("data", info)
		c.RenderData("count", count)
		c.Success(c.GetRenderData())
	}
	return f
}
