/**
 * Created by angelina on 2017/9/5.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package easyweb_cms

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vannnnish/easyweb"
	"github.com/vannnnish/easyweb_cms/conf"
	"github.com/vannnnish/easyweb_cms/model"
	"github.com/vannnnish/easyweb_cms/route"
	"github.com/vannnnish/yeego"
	"github.com/vannnnish/yeego/yeecrypto"
	"github.com/vannnnish/yeego/yeestrconv"
	"io"
	"log"
)

var defaultDB *gorm.DB

type DbConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	DbName   string
}

// 创建数数据库
func CreateDb(conf DbConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName, conf.Password, conf.Host, conf.Port)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARSET utf8mb4 COLLATE utf8mb4_general_ci", conf.DbName)
	err = db.Exec(sql).Error
	if err != nil {
		return err
	}
	db.CreateTable()
	return nil
}

// 初始化gorm的DB
func InitDefaultDB(conf DbConfig) error {
	CreateDb(conf)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName, conf.Password, conf.Host, conf.Port, conf.DbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(2000)
	model.InitDB(db)
	defaultDB = db
	return nil
}

// 是否输出数据库语句
func SetDbLog(b bool) {
	defaultDB.LogMode(b)
}

func SetDbLogWriter(w io.Writer) {
	defaultDB.SetLogger(log.New(w, "\r\n", 0))
}

func GetDB() *gorm.DB {
	return defaultDB
}

// 创建配置的数据库表
func InitDefaultTable() {
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.AdminLog{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.AdminPrivilege{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.AdminRole{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.AdminRolePrivilege{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.AdminUserToken{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.AdminUser{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.Article{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.Category{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.Kvdb{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.Model{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.SinglePage{})
	defaultDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.Slide{})
}

func InitDefaultRoute(web *easyweb.EasyWeb, dispatchMap map[string]map[int]easyweb.HandlerFunc) {
	route.InitRoute(web, dispatchMap)
}

func initCategory(category *model.Category, privilege string) {
	var err error
	err = defaultDB.FirstOrCreate(category).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
	}
	err = model.AdminPrivilege{}.CreateOrUpdate(category.Id, privilege)
	if err != nil {
		easyweb.Logger.Error(err.Error())
	}
}

func InitDefaultCategory() {
	// 内容管理
	category := &model.Category{
		Id:       conf.ContentCateId,
		Name:     "内容管理",
		ParentId: conf.AdminTopCateId,
		ModelId:  conf.DirModelId,
		Sort:     3,
		Contain:  yeego.Config.GetString("app.Contains"),
	}
	initCategory(category, "profile")
	// 栏目管理
	category = &model.Category{
		Id:       conf.CategoryCateId,
		Name:     "栏目管理",
		ParentId: conf.AdminTopCateId,
		ModelId:  conf.CategoryModelId,
		Sort:     1,
	}
	initCategory(category, "create,delete,profile,sort,update")
	// 管理员管理
	category = &model.Category{
		Id:       conf.AdminCateId,
		Name:     "管理员管理",
		ParentId: conf.AdminTopCateId,
		ModelId:  conf.DirModelId,
		Sort:     2,
	}
	initCategory(category, "profile")
	// 角色管理
	category = &model.Category{
		Id:       conf.AdminRoleCateId,
		Name:     "角色",
		ParentId: conf.AdminCateId,
		ModelId:  conf.AdminRoleModelId,
		Sort:     3,
	}
	initCategory(category, "create,delete,profile,sort,update")
	// 用户管理
	category = &model.Category{
		Id:       conf.AdminUserCateId,
		Name:     "后台用户",
		ParentId: conf.AdminCateId,
		ModelId:  conf.AdminUserModelId,
		Sort:     2,
	}
	initCategory(category, "create,delete,profile,sort,update")
	// 日志管理
	category = &model.Category{
		Id:       conf.AdminLogCateId,
		Name:     "操作日志",
		ParentId: conf.AdminCateId,
		ModelId:  conf.AdminUserLogModelId,
		Sort:     1,
	}
	initCategory(category, "profile")
}

func getConfigModel() []map[string]string {
	models := make(map[string][]map[string]string)
	_, err := toml.DecodeFile(conf.ModelConfigPath, &models)
	if err != nil {
		panic(err)
	}
	return models["model"]
}

func InitDefaultModel() {
	models := getConfigModel()
	for _, v := range models {
		if v["IsNeed"] == "1" {
			var isShow bool
			if v["IsShow"] == "1" {
				isShow = true
			} else {
				isShow = false
			}
			m := model.Model{
				Id:          yeestrconv.AtoIDefault0(v["Id"]),
				Name:        v["Name"],
				IsShow:      isShow,
				DbTableName: v["TableName"],
				Actions:     v["Actions"],
			}
			err := model.Model{}.Create(m)
			if err != nil {
				easyweb.Logger.Error(err.Error())
			}
		}
	}
}

func InitDefaultAdminUser() {
	user1 := &model.AdminUser{
		Id:       1,
		Account:  "yeeyun_root",
		RoleId:   -1,
		Password: yeecrypto.Sha256Hex([]byte(yeecrypto.Sha256Hex([]byte("cd32d5e86e")))),
	}
	defaultDB.FirstOrCreate(user1)
	user2 := &model.AdminUser{
		Id:       2,
		Account:  "yeeyun",
		RoleId:   -1,
		Password: yeecrypto.Sha256Hex([]byte(yeecrypto.Sha256Hex([]byte("cd32d5e86e")))),
	}
	defaultDB.FirstOrCreate(user2)
}

func InitAllDefault() {
	InitDefaultTable()
	InitDefaultCategory()
	InitDefaultModel()
	InitDefaultAdminUser()
}
