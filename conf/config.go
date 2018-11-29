/**
 * Created by angelina on 2017/9/8.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package conf

const (
	// 超级管理员的角色id
	SuperAdminRoleId int = -1
	// 模型缓存文件目录
	CacheFilePath = "cache/model/"
	// 栏目缓存文件目录
	CategoryFilePath = "cache/category/"
	// 背景图片缓存地址
	BgImageCachePath = "cache/bgimage/bgimage.cache"
	// 文件上传根目录
	UploadPath = "data/upload/"
	// 文件上传路径
	UploadFilePath = UploadPath + "file/"
	// 图片上传路径
	UploadImgPath = UploadPath + "image/"
	// 视频上传路径
	UploadVideoPath = UploadPath + "video/origin/"
	// 视频转码存储路径
	UploadVideoTranscodePath = UploadPath + "video/transcode/"
	// ueditor文件存储路径
	UeditorUploadPath = UploadPath + "ueditor/"
	// 模型配置文件路径
	ModelConfigPath = "conf/model.toml"
)

// 模型id
const (
	// 目录模型id
	DirModelId = 1
	// 管理员角色模型id
	AdminRoleModelId = 2
	// 管理员日志模型id
	AdminUserLogModelId = 3
	// 管理员模型id
	AdminUserModelId = 4
	// 栏目模型id
	CategoryModelId = 5
	// 文章模型id
	ArticleModelId = 6
	// 单页面模型id
	SinglePageModelId = 7
	// 轮换图模型id
	SlideModelId = 8
	// 下载模型id
	DownloadModelId = 9
)

// 栏目id
const (
	// 后台栏目顶级栏目id
	AdminTopCateId = 0
	// 内容管理
	ContentCateId = 1
	// 栏目管理
	CategoryCateId = 2
	// 管理员管理
	AdminCateId = 3
	// 角色管理
	AdminRoleCateId = 4
	// 用户管理
	AdminUserCateId = 5
	// 日志管理
	AdminLogCateId = 6
)
