/**
 * Created by WillkYang on 2018/10/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */
package api

import (
	"github.com/spf13/cast"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
	"x-market_lib/util/json_tool"
)

var CommonApiCtx CommonApi

type CommonApi struct {
}

func (CommonApi) Sort() easyweb.HandlerFunc {
	return func(c *easyweb.Context) {
		modelId := c.PathParam("modelId")
		sort := c.Param("sort").GetString()
		var set map[int]int
		err := json_tool.Json().UnmarshalFromString(sort, &set)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		for k, v := range set {
			err := model.DoSort(cast.ToInt(modelId), k, v)
			if err != nil {
				c.FailWithDefaultCode(err.Error())
				return
			}
		}
		c.SuccessWithMsg("排序成功")
	}
}
