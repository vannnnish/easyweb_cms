/**
 * Created by WillkYang on 2018/10/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */
package api

import (
	"encoding/json"
	"github.com/spf13/cast"
	"gitlab.yeeyuntech.com/yee/easyweb"
	"gitlab.yeeyuntech.com/yee/easyweb_cms/model"
)

var CommonApiCtx CommonApi

type CommonApi struct {
}

func (CommonApi) Sort() easyweb.HandlerFunc {
	return func(c *easyweb.Context) {
		modelId := c.PathParam("modelId")
		sort := c.Param("sort").GetString()
		var set map[int]int
		err := json.Unmarshal([]byte(sort), &set)
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
