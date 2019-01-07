/**
 * Created by angelina on 2017/10/21.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package api

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/easyweb_cms/conf"
	"github.com/vannnnish/yeego/yeeCrypto"
	"github.com/vannnnish/yeego/yeeFile"
	"github.com/vannnnish/yeego/yeeImage"
	"github.com/vannnnish/yeego/yeeStrconv"
	"github.com/vannnnish/yeego/yeeStrings"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Upload_Api struct {
}

func (Upload_Api) UploadImage(saveDir ...string) easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		savePath := conf.UploadImgPath
		if len(saveDir) > 0 {
			savePath = saveDir[0]
		}
		returnPath := conf.UploadImgPath
		if len(saveDir) > 1 {
			returnPath = saveDir[1]
		}
		info := SimpleUpload(c.Request(), savePath, returnPath, ImageExtList)
		if info.Err != nil {
			c.FailWithDefaultCode(info.Err.Error())
			return
		}
		thumbnails := c.Param("thumb").GetString()
		if thumbnails != "" {
			var arr [][]int
			err := json.Unmarshal([]byte(thumbnails), &arr)
			if err != nil {
				easyweb.Logger.Error(err.Error())
			} else {
				for _, v := range arr {
					if len(v) != 2 {
						continue
					}
					pathArr := strings.Split(info.Url, ".")
					ext := "." + pathArr[len(pathArr)-1]
					pathPrefix := info.Url[:len(info.Url)-len(ext)]
					thumbPath := "." + pathPrefix + "_" + yeeStrconv.FormatInt(v[0]) + "x" + yeeStrconv.FormatInt(v[1]) + ext
					yeeImage.ResizeImage("."+info.Url, thumbPath, uint(v[0]), uint(v[1]))
				}
			}
		}
		c.Success(info)
	}
	return f
}

func (Upload_Api) UploadFile(saveDir ...string) easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		savePath := conf.UploadFilePath
		if len(saveDir) > 0 {
			savePath = saveDir[0]
		}
		returnPath := conf.UploadFilePath
		if len(saveDir) > 1 {
			returnPath = saveDir[1]
		}
		info := SimpleUpload(c.Request(), savePath, returnPath, ImageExtList)
		if info.Err != nil {
			c.FailWithDefaultCode(info.Err.Error())
			return
		}
		c.Success(info)
	}
	return f
}

func (Upload_Api) UploadVideo(saveDir ...string) easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		info := SimpleUpload(c.Request(), conf.UploadVideoPath, "", ImageExtList)
		if info.Err != nil {
			c.FailWithDefaultCode(info.Err.Error())
			return
		}
		transcode := c.Param("transcode").SetDefault("false").GetBool()
		if transcode {
			dst := genRealPathWithDir(conf.UploadVideoTranscodePath) + info.FileName
			go transCodeToMp4(info.Url[1:], dst)
			info.Url = "/" + dst
		}
		c.Success(info)
	}
	return f
}

type UploadReturnInfo struct {
	Url        string // 地址
	FileName   string // 新文件名称
	Ext        string // 后缀
	OriginName string // 原始文件名称
	Err        error  // 错误
}

var (
	ImageExtList  = []string{".jpg", ".png", ".gif", ".jpeg"}
	VideoExtLIst  = []string{".flv", ".swf", ".mkv", ".avi", ".rmvb", ".mov", ".wmv", ".mp4", ".mp3"}
	AttachExtList = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp",
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg", ".ogg",
		".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
		".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf",
		".txt", ".md", ".xml", ".exe"}
	ContentTypeExt = map[string]string{
		"application/vnd.ms-excel":                                                  ".xls",
		"application/vnd.ms-powerpoint":                                             ".ppt",
		"application/vnd.ms-word":                                                   ".doc",
		"application/vnd.oasis.chart":                                               ".odc",
		"application/vnd.oasis.database":                                            ".odb",
		"application/vnd.oasis.formula":                                             ".odf",
		"application/vnd.oasis.image":                                               ".odi",
		"application/vnd.oasis.opendocument.graphics":                               ".odg",
		"application/vnd.oasis.opendocument.graphics-template":                      ".otg",
		"application/vnd.oasis.opendocument.presentation":                           ".odp",
		"application/vnd.oasis.opendocument.presentation-template":                  ".otp",
		"application/vnd.oasis.opendocument.text":                                   ".odt",
		"application/vnd.oasis.opendocument.text-master":                            ".odm",
		"application/vnd.oasis.opendocument.text-template":                          ".ott",
		"application/vnd.oasis.opendocument.text-web":                               ".oth",
		"application/vnd.oasis.spreadsheet":                                         ".ods",
		"application/vnd.oasis.spreadsheet-template":                                ".ots",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
		"application/vnd.openxmlformats-officedocument.presentationml.slideshow":    ".ppsx",
		"application/vnd.openxmlformats-officedocument.presentationml.template":     ".potx",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.template":      ".xltx",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   ".docx",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.template":   ".dotx",
		"application/x-ole-storage":                                                 ".msg",
		"image/gif":                                                                 ".gif",
		"image/jpeg":                                                                ".jpg",
		"image/png":                                                                 ".png",
		"text/plain":                                                                ".txt",
	}
)

// SimpleUpload
// 简单的上传器
func SimpleUpload(request *http.Request, uploadDir, returnPath string, allowExt []string) (info UploadReturnInfo) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("simple upload error :[%v] \n", err)
		}
	}()
	if request.Method != "POST" {
		info.Err = errors.New("必须是POST请求")
		checkErr(errors.New("必须是POST请求"))
	}
	// 原始文件名称
	var inputFileName string
	// 文件后缀
	var inputFileExt string
	// md5后的文件名称
	var outputFileName string
	if returnPath == "" {
		returnPath = uploadDir
	}
	if uploadDir == "" {
		uploadDir = genRealPathWithDir(conf.UploadPath)
		returnPath = uploadDir
	} else {
		uploadDir = genRealPathWithDir(uploadDir)
		returnPath = genRealPathWithDirWithoutMkdir(returnPath)
	}
	os.MkdirAll(uploadDir, os.ModePerm)
	// 此处不忽略file是因为FormFile已经创建了file，如果忽略的话，不对file释放，会造成内存泄漏???
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		info.Err = err
		checkErr(err)
	}
	file.Close()
	inputFileName = fileHeader.Filename
	inputFileExt = filepath.Ext(inputFileName)
	if inputFileExt == "" {
		mime := fileHeader.Header.Get("Content-Type")
		if mime != "" {
			if ext := ContentTypeExt[mime]; ext != "" {
				inputFileName = inputFileName + ext
				inputFileExt = filepath.Ext(inputFileName)
			}
		}
	}
	if !yeeStrings.IsInSlice(allowExt, strings.ToLower(inputFileExt)) {
		info.Err = err
		checkErr(err)
		return UploadReturnInfo{Err: errors.New("文件后缀不被允许")}
	}
	inputFile, err := fileHeader.Open()
	if err != nil {
		info.Err = err
		checkErr(err)
	}
	defer inputFile.Close()
	outputFileName = yeeCrypto.Md5Hex([]byte(inputFileName))
	// 根据是否分片上传决定上传策略
	// 分片上传路径目录
	var chunksUploadDir string
	// 文件的存储地址
	var outputFilePath string
	if len(request.FormValue("chunks")) > 0 {
		chunksUploadDir = filepath.Join(uploadDir, outputFileName) + "/"
		os.MkdirAll(chunksUploadDir, os.ModePerm)
		chunkFileName := fmt.Sprintf("%s-%s", outputFileName, request.FormValue("chunk"))
		outputFilePath = filepath.Join(chunksUploadDir, chunkFileName)
	} else {
		md5h := md5.New()
		_, err = io.Copy(md5h, inputFile)
		if err != nil {
			info.Err = err
			checkErr(err)
		}
		outputFileName = fmt.Sprintf("%x%s", md5h.Sum(nil), inputFileExt)
		outputFilePath = filepath.Join(uploadDir, outputFileName)
		returnPath = filepath.Join(returnPath, outputFileName)
	}
	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		info.Err = err
		checkErr(err)
	}
	defer outputFile.Close()
	inputFile.Seek(0, 0)
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		info.Err = err
		checkErr(err)
	}
	// 根据上传策略修改文件分片计数,合并
	if len(request.FormValue("chunks")) > 0 {
		chunks := yeeStrconv.AtoIDefault0(request.FormValue("chunks"))
		chunkCount := fileCount(filepath.Join(chunksUploadDir, yeeStrconv.FormatInt(chunks)))
		if chunks == chunkCount {
			outputFilePath, err = mergeChunkFile(chunksUploadDir, uploadDir)
			if err != nil {
				info.Err = err
				checkErr(err)
			}
			fileNames, _ := com.StatDir(chunksUploadDir)
			mergeFileName := filepath.Base(chunksUploadDir)
			returnPath = uploadDir + mergeFileName + filepath.Ext(fileNames[0])
		}
	}
	if returnPath[0] != '/' {
		returnPath = "/" + returnPath
	}
	return UploadReturnInfo{
		Url:        returnPath,
		FileName:   outputFileName,
		Ext:        inputFileExt,
		OriginName: inputFileName,
		Err:        nil,
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GenRealPathWithDir
// 在目录后加上年月
func genRealPathWithDir(dir string) string {
	now := time.Now()
	newDir := fmt.Sprintf("%s/%d/%d/%d/", dir, now.Year(), now.Month(), now.Day())
	if !yeeFile.FileExists(newDir) {
		os.MkdirAll(newDir, os.ModePerm)
	}
	return newDir
}

func genRealPathWithDirWithoutMkdir(dir string) string {
	now := time.Now()
	newDir := fmt.Sprintf("%s/%d/%d/%d/", dir, now.Year(), now.Month(), now.Day())
	return newDir
}

var fileCountLock sync.Mutex

// fileCount
// 文件计数，统计目前有多少个分片已经上传成功
func fileCount(path string) int {
	fileCountLock.Lock()
	defer fileCountLock.Unlock()
	if yeeFile.FileExists(path) {
		countStr, err := yeeFile.GetString(path)
		if err != nil {
			count := 1
			yeeFile.SetString(path, yeeStrconv.FormatInt(count))
			return count
		}
		count := yeeStrconv.AtoIDefault0(countStr) + 1
		yeeFile.SetString(path, yeeStrconv.FormatInt(count))
		return count
	}
	count := 1
	yeeFile.SetString(path, yeeStrconv.FormatInt(count))
	return count
}

// mergeChunkFile
// 合并分片文件
func mergeChunkFile(inputDir, outputDir string) (string, error) {
	fileNames, _ := com.StatDir(inputDir)
	mergeFileName := filepath.Base(inputDir)
	newFilePath := outputDir + mergeFileName + filepath.Ext(fileNames[0])
	os.RemoveAll(newFilePath)
	outputFile, err := os.OpenFile(newFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	if len(fileNames) > 1 {
		for i := 0; i < len(fileNames)-1; i++ {
			if err := fileCopy(fmt.Sprintf("%s%s-%d%s", inputDir, filepath.Base(inputDir), i,
				filepath.Ext(fileNames[0])), outputFile); err != nil {
				return "", err
			}
		}
	}
	os.RemoveAll(inputDir)
	return newFilePath, nil
}

// fileCopy
// 文件拷贝
func fileCopy(fileName string, outputFile *os.File) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(outputFile, file)
	if err != nil {
		return err
	}
	return nil
}

// transCodeToMp4
// 视频转码为MP4
func transCodeToMp4(src, dst string) error {
	// 视频转码
	// 使用qt-faststart来把meta信息移到文件头部
	command := fmt.Sprintf("avconv -y -i %s -vcodec h264 -strict -2 %s;"+
		"qt-faststart %s %s.faststart", src, dst, dst, dst)
	cmd := exec.Command("/bin/sh", "-c", command)
	err := cmd.Start()
	if err != nil {
		fmt.Println(fmt.Sprintf("transCodeToMp4 error :[%s]", err.Error()))
		return err
	}
	_, err = cmd.Process.Wait()
	if err != nil {
		fmt.Println(fmt.Sprintf("transCodeToMp4 error :[%s]", err.Error()))
		return err
	}
	// dst不存在，表示转码失败，则直接拷贝原始视频到新地址
	if !yeeFile.FileExists(dst) {
		if err := yeeFile.Copy(src, dst); err != nil {
			fmt.Println(fmt.Sprintf("transCodeToMp4 error :[%s]", err.Error()))
			return err
		}
	}
	// 存在dst.faststart 说明qt-faststart执行成功，则换名称
	if yeeFile.FileExists(dst + ".faststart") {
		// 先copy一份旧的视频,拷贝失败，直接返回，不管了
		if err := yeeFile.Copy(dst, dst+".old"); err != nil {
			fmt.Println(fmt.Sprintf("transCodeToMp4 error :[%s]", err.Error()))
			os.RemoveAll(dst + ".faststart")
			return nil
		}
		// 将qt-start的改名
		if err := os.Rename(dst+".faststart", dst); err != nil {
			fmt.Println(fmt.Sprintf("transCodeToMp4 error :[%s]", err.Error()))
			// 如果失败，那么就再转回去
			yeeFile.Copy(dst+".old", dst)
		}
		os.RemoveAll(dst + ".old")
	}
	if !yeeFile.FileExists(dst) {
		if err := yeeFile.Copy(src, dst); err != nil {
			fmt.Println(fmt.Sprintf("transCodeToMp4 error :[%s]", err.Error()))
			return err
		}
	}
	return nil
}
