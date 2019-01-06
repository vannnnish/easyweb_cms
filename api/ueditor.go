/**
 * Created by angelina on 2017/10/21.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"easyweb_cms/conf"
	"encoding/json"
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/yeego"
	"sync"
)

type Ueditor_Api struct {
}

func (Ueditor_Api) UE() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		action := c.Param("action").GetString()
		switch action {
		case "config":
			ueditorConfLock.Lock()
			config := getFrontConfig()
			callback := c.Param("callback").GetString()
			configJson, err := json.Marshal(config)
			ueditorConfLock.Unlock()
			if err != nil {
				c.String(200, err.Error())
				return
			}
			c.String(200, callback+"("+string(configJson)+")")
		case "uploadimage":
			path := yeego.Config.GetString("upload.UeditorImgPath")
			if path == "" {
				path = conf.UeditorUploadPath
			}
			returnPath := yeego.Config.GetString("upload.UeditorImgReturnPath")
			if returnPath == "" {
				returnPath = conf.UeditorUploadPath
			}
			info := SimpleUpload(c.Request(), path, returnPath, ImageExtList)
			rInfo := map[string]interface{}{
				"state":    "SUCCESS",
				"title":    info.FileName,
				"original": info.FileName,
				"url":      info.Url,
			}
			c.JSON(200, rInfo)
		case "uploadfile":
			path := yeego.Config.GetString("upload.UeditorFilePath")
			if path == "" {
				path = conf.UeditorUploadPath
			}
			returnPath := yeego.Config.GetString("upload.UeditorFileReturnPath")
			if returnPath == "" {
				returnPath = conf.UeditorUploadPath
			}
			info := SimpleUpload(c.Request(), path, returnPath, AttachExtList)
			rInfo := map[string]interface{}{
				"state":    "SUCCESS",
				"title":    info.FileName,
				"original": info.FileName,
				"url":      info.Url,
			}
			c.JSON(200, rInfo)
		case "uploadvideo":
			path := yeego.Config.GetString("upload.UeditorVideoPath")
			if path == "" {
				path = conf.UeditorUploadPath
			}
			returnPath := yeego.Config.GetString("upload.UeditorVideoReturnPath")
			if returnPath == "" {
				returnPath = conf.UeditorUploadPath
			}
			info := SimpleUpload(c.Request(), path, returnPath, VideoExtLIst)
			rInfo := map[string]interface{}{
				"state":    "SUCCESS",
				"title":    info.FileName,
				"original": info.FileName,
				"url":      info.Url,
			}
			c.JSON(200, rInfo)
		default:
			c.JSON(200, map[string]string{"err": "action not found"})
		}
	}
	return f
}

var UeditorConfig = map[string]interface{}{
	"imageFieldName":      "file",
	"imageActionName":     "uploadimage",
	"imageMaxSize":        2048000,
	"imageAllowFiles":     []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"},
	"imageCompressEnable": true,
	"imageCompressBorder": 1600,
	"imageInsertAlign":    "none",
	"imageUrlPrefix":      "",
	"imagePathFormat":     "/ueditor/php/upload/image/{yyyy}{mm}{dd}/{time}{rand:6}",

	"scrawlActionName":  "uploadscrawl",
	"scrawlFieldName":   "file",
	"scrawlPathFormat":  "/ueditor/php/upload/image/{yyyy}{mm}{dd}/{time}{rand:6}",
	"scrawlMaxSize":     2048000,
	"scrawlUrlPrefix":   "",
	"scrawlInsertAlign": "none",

	"snapscreenActionName":  "uploadimage",
	"snapscreenPathFormat":  "/ueditor/php/upload/image/{yyyy}{mm}{dd}/{time}{rand:6}",
	"snapscreenUrlPrefix":   "",
	"snapscreenInsertAlign": "none",

	"catcherLocalDomain": []string{"127.0.0.1", "localhost", "img.baidu.com"},
	"catcherActionName":  "catchimage",
	"catcherFieldName":   "source",
	"catcherPathFormat":  "/ueditor/php/upload/image/{yyyy}{mm}{dd}/{time}{rand:6}",
	"catcherUrlPrefix":   "",
	"catcherMaxSize":     2048000,
	"catcherAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"},

	"videoActionName": "uploadvideo",
	"videoFieldName":  "file",
	"videoPathFormat": "/ueditor/php/upload/video/{yyyy}{mm}{dd}/{time}{rand:6}",
	"videoUrlPrefix":  "",
	"videoMaxSize":    102400000,
	"videoAllowFiles": []string{
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
		".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid"},

	"fileActionName": "uploadfile",
	"fileFieldName":  "file",
	"filePathFormat": "/ueditor/php/upload/file/{yyyy}{mm}{dd}/{time}{rand:6}",
	"fileUrlPrefix":  "",
	"fileMaxSize":    51200000,
	"fileAllowFiles": []string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp",
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
		".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
		".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		".exe"},
	"imageManagerActionName":  "listimage",
	"imageManagerListPath":    "/ueditor/php/upload/image/",
	"imageManagerListSize":    20,
	"imageManagerUrlPrefix":   "",
	"imageManagerInsertAlign": "none",
	"imageManagerAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"},
	"fileManagerActionName":   "listfile",
	"fileManagerListPath":     "/ueditor/php/upload/file/",
	"fileManagerUrlPrefix":    "",
	"fileManagerListSize":     20,
	"fileManagerAllowFiles": []string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp",
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
		".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
		".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		".exe"},
}

var ueditorConfLock sync.RWMutex

func getFrontConfig() map[string]interface{} {
	var config map[string]interface{} = UeditorConfig
	config["imagePathFormat"] = ImageExtList
	config["catcherAllowFiles"] = ImageExtList
	config["imageManagerAllowFiles"] = ImageExtList
	config["fileAllowFiles"] = AttachExtList
	config["fileManagerAllowFiles"] = AttachExtList
	return config
}
