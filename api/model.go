/**
 * Created by angelina on 2017/9/13.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
	"gitlab.yeeyuntech.com/yee/easyweb"
)

type Model_Api struct {
	ModelModel model.Model
}

// 模型列表信息
func (m Model_Api) List() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		models := m.ModelModel.SelectAllWithCache()
		c.Success(models)
	}
	return f
}

// 单个模型信息
func (m Model_Api) Profile() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		modelId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		model := &model.Model{Id: modelId}
		m.ModelModel.SelectOneModelWithCache(model)
		c.Success(model)
	}
	return f
}
