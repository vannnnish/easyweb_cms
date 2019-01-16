/**
 * Created by angelina on 2017/9/8.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/easyweb_cms/model"
	"github.com/vannnnish/yeego/yeestrconv"
)

type Article_Api struct {
	ArticleModel model.Article
}

// 获取文章列表
func (article Article_Api) List() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.PathParam("cate_id")
		keyword := c.Param("keyword").GetString()
		source := c.Param("source").GetString()
		author := c.Param("author").GetString()
		isPublish := c.Param("is_publish").SetDefaultInt(2).GetInt()
		startTime := c.Param("start_time").GetString()
		endTime := c.Param("end_time").GetString()
		pageSize := c.Param("pagesize").SetDefaultInt(DefaultPageSize).GetInt()
		page := c.Param("page").SetDefaultInt(DefaultPage).GetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		count := article.ArticleModel.SelectAllCount(yeestrconv.AtoIDefault0(cateId), keyword, source, author, isPublish, startTime, endTime)
		data := article.ArticleModel.SelectAll(yeestrconv.AtoIDefault0(cateId), keyword, source, author, isPublish, startTime, endTime,
			pageSize, (page-1)*pageSize)
		c.RenderData("data", data)
		c.RenderData("count", count)
		c.Success(c.GetRenderData())
	}
	return f
}

// 文章详情
func (article Article_Api) Profile() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		articleId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		info := &model.Article{Id: articleId}
		article.ArticleModel.SelectOne(info)
		c.Success(info)
	}
	return f
}

// 创建文章
func (article Article_Api) Create() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		cateId := c.PathParam("cate_id")
		title := c.Param("title").MustGetString()
		thumb := c.Param("thumb").GetString()
		source := c.Param("source").GetString()
		author := c.Param("author").GetString()
		picAuthor := c.Param("pic_author").GetString()
		desc := c.Param("description").GetString()
		content := c.Param("content").MustGetString()
		isPublish := c.Param("is_publish").GetBool()
		updateTime := c.Param("update_time").GetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		newArticle := model.Article{
			CateId:      yeestrconv.AtoIDefault0(cateId),
			Title:       title,
			Thumb:       thumb,
			Source:      source,
			Author:      author,
			PicAuthor:   picAuthor,
			Description: desc,
			Content:     content,
			IsPublish:   isPublish,
			UpdateTime:  updateTime,
		}
		err := article.ArticleModel.Create(newArticle)
		if err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("添加文章成功")
	}
	return f
}

// 编辑文章
func (article Article_Api) Update() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		articleId := c.Param("id").MustGetInt()
		title := c.Param("title").MustGetString()
		thumb := c.Param("thumb").GetString()
		source := c.Param("source").GetString()
		author := c.Param("author").GetString()
		picAuthor := c.Param("pic_author").GetString()
		desc := c.Param("description").GetString()
		content := c.Param("content").MustGetString()
		isPublish := c.Param("is_publish").GetBool()
		updateTime := c.Param("update_time").GetString()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		newArticle := model.Article{
			Id:          articleId,
			Title:       title,
			Thumb:       thumb,
			Source:      source,
			Author:      author,
			PicAuthor:   picAuthor,
			Description: desc,
			Content:     content,
			IsPublish:   isPublish,
			UpdateTime:  updateTime,
		}
		if err := article.ArticleModel.Update(newArticle); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("编辑文章成功")
	}
	return f
}

// 排序文章
func (article Article_Api) Sort() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		postData := c.Request().PostForm
		for k, v := range postData {
			err := article.ArticleModel.DoSort(yeestrconv.AtoIDefault0(k), yeestrconv.AtoIDefault0(v[0]))
			if err != nil {
				c.FailWithDefaultCode(err.Error())
				return
			}
		}
		c.SuccessWithMsg("排序成功")
	}
	return f
}

// 发布文章
func (article Article_Api) Publish() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		articleId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := article.ArticleModel.Publish(articleId); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("发布文章成功")
	}
	return f
}

// 取消发布文章
func (article Article_Api) UnPublish() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		articleId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := article.ArticleModel.UnPublish(articleId); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("取消发布文章成功")
	}
	return f
}

// 删除文章
func (article Article_Api) Delete() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		articleId := c.Param("id").MustGetInt()
		if err := c.GetError(); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		if err := article.ArticleModel.Delete(articleId); err != nil {
			c.FailWithDefaultCode(err.Error())
			return
		}
		c.SuccessWithMsg("删除文章成功")
	}
	return f
}
