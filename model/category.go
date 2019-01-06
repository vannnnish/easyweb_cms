/**
 * Created by angelina on 2017/9/13.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package model

import (
	"encoding/json"
	"errors"
	"github.com/vannnnish/yeego/yeeCache"
	"github.com/vannnnish/yeego/yeeSql"
	"github.com/vannnnish/yeego/yeeStrconv"
	"github.com/vannnnish/yeego/yeeStrings"
	"github.com/vannnnish/yeego/yeeTransform"
	"github.com/vannnnish/easyweb"
	"easyweb_cms/conf"
	"time"
)

// 栏目信息
type Category struct {
	Id       int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"` //
	Name     string `gorm:"size:50;not null" json:"name"`         // 栏目名称
	ParentId int    `gorm:"not null" json:"parent_id"`            // 父亲栏目id
	ModelId  int    `gorm:"not null" json:"model_id"`             // 所属模型id
	Sort     int    `gorm:"not null" json:"sort"`                 // 排序
	Contain  string `gorm:"not null;size:100" json:"contain"`     // 如果是一个目录，它下面可以包含哪些model类型
}

func (Category) TableName() string {
	return "yeecms_category"
}

// 查找栏目
func (Category) SelectOne(category *Category) error {
	if err := defaultDB.First(category, category.Id).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 通过缓存查找栏目
func (Category) SelectOneWithCache(category *Category) error {
	cacheFileName := conf.CategoryFilePath + yeeStrconv.FormatInt(category.Id) + ".cache"
	b, err := yeeCache.FileTtlCache(cacheFileName, func() (b []byte, ttl time.Duration, err error) {
		err = Category{}.SelectOne(category)
		if err != nil {
			return nil, time.Second, err
		}
		b, err = json.Marshal(category)
		if err != nil {
			return nil, time.Second, err
		}
		ttl = time.Hour * 24
		return
	})
	if err != nil {
		if err := (Category{}).SelectOne(category); err != nil {
			return err

		}
		return nil
	}
	if err := json.Unmarshal(b, category); err != nil {
		easyweb.Logger.Error(err.Error())
		return SelectError
	}
	return nil
}

// 获取栏目的Actions,即这个栏目拥有哪些操作权限(create,update,profile...等等)
func (Category) GetRoleActions(cateId int) string {
	sql := `
			SELECT
				group_concat(action)                           AS actions,
				group_concat(yeecms_privilege.id ORDER BY yeecms_privilege.id ASC) AS action_ids
			FROM
				(
				  SELECT *
				  FROM yeecms_category
				  WHERE id = ?
				) category
			LEFT JOIN
				yeecms_privilege
			ON
				category.id = yeecms_privilege.cate_id
			GROUP BY
				category.id;
	`
	rows, err := defaultDB.Raw(sql, cateId).Rows()
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return ""
	}
	info, err := yeeSql.RowsToMapSliceFirst(rows)
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return ""
	}
	return info["actions"]
}

// 查找全部的栏目(以及某个角色拥有的栏目权限)
func (Category) SelectAll(parentId, roleId int, isRecursion bool) []map[string]interface{} {
	var categories []map[string]string
	if roleId == conf.SuperAdminRoleId {
		sql := `
			SELECT
			  category.*,
			  privilege.actions AS role_actions
			FROM
			  (SELECT
				 category.*,
				 group_concat(action)                           AS actions,
				 group_concat(yeecms_privilege.id ORDER BY yeecms_privilege.id ASC) AS action_ids
			   FROM
				 (SELECT *
				  FROM yeecms_category
				  WHERE parent_id = ?) category
				 LEFT JOIN yeecms_privilege
				   ON category.id = yeecms_privilege.cate_id
			   GROUP BY category.id) category
			  LEFT JOIN
			  (SELECT
				 cate_id,
				 group_concat(action) AS actions
			   FROM yeecms_privilege
			   GROUP BY cate_id) privilege
				ON category.id = privilege.cate_id
			ORDER BY sort DESC, category.id ASC
		`
		rows, err := defaultDB.Raw(sql, parentId).Rows()
		if err != nil {
			easyweb.Logger.Error(err)
			return nil
		}
		categories = yeeSql.RowsToMapSlice(rows)
	} else {
		sql := `
			SELECT
			  category.*,
			  privilege.actions AS role_actions
			FROM
			  (SELECT
				 category.*,
				 group_concat(action)                           AS actions,
				 group_concat(yeecms_privilege.id ORDER BY yeecms_privilege.id ASC) AS action_ids
			   FROM
				 (SELECT *
				  FROM yeecms_category
				  WHERE parent_id = ?) category
				 LEFT JOIN yeecms_privilege
				   ON category.id = yeecms_privilege.cate_id
			   GROUP BY category.id) category
			  LEFT JOIN
			  (SELECT
				 cate_id,
				 group_concat(action) AS actions
			   FROM yeecms_privilege, yeecms_admin_role_privilege
			   WHERE yeecms_privilege.id = yeecms_admin_role_privilege.privilege_id
			   AND role_privilege.role_id = ?
			   GROUP BY cate_id) privilege
				ON category.id = privilege.cate_id
			ORDER BY sort DESC, category.id ASC;
		`
		rows, _ := defaultDB.Raw(sql, parentId, roleId).Rows()
		categories = yeeSql.RowsToMapSlice(rows)
	}
	for _, category := range categories {
		pId := yeeStrconv.AtoIDefault0(category["parent_id"])
		if pId != 0 {
			parentCate := &Category{Id: pId}
			Category{}.SelectOneWithCache(parentCate)
			category["parent_cate_name"] = parentCate.Name
		} else {
			category["parent_cate_name"] = ""
		}
	}
	data := yeeTransform.MapSliceStringToInterface(categories)
	if isRecursion {
		for _, v := range data {
			children := Category{}.SelectAll(yeeStrconv.AtoIDefault0(v["id"].(string)), roleId, isRecursion)
			if len(children) == 0 {
				v["children"] = nil
			} else {
				v["children"] = children
			}
		}
	}
	return data
}

// 根据缓存查找全部的栏目数据
func (Category) SelectAllWithCache(parentId, roleId int, isRecursion bool) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	cacheFileName := conf.CategoryFilePath + "all.cache"
	b, err := yeeCache.FileTtlCache(cacheFileName, func() (b []byte, ttl time.Duration, err error) {
		data := Category{}.SelectAll(parentId, roleId, isRecursion)
		b, err = json.Marshal(&data)
		if err != nil {
			return nil, time.Second, err
		}
		ttl = time.Hour * 24
		return
	})
	if err != nil {
		return Category{}.SelectAll(parentId, roleId, isRecursion)
	}
	if err := json.Unmarshal(b, &data); err != nil {
		easyweb.Logger.Error(err.Error())
		return nil
	}
	return data
}

// 创建栏目
func (Category) Create(category Category) (int, error) {
	var cateId int
	parentCategory := &Category{Id: category.ParentId}
	if err := (Category{}).SelectOne(parentCategory); err != nil {
		return cateId, err
	}
	// 父级栏目的模型必须是 目录  或者  栏目管理
	if !(parentCategory.ModelId == conf.DirModelId || parentCategory.ModelId == conf.CategoryModelId) {
		return cateId, errors.New("只允许在目录下创建栏目")
	}
	// 判断模型是否被允许
	if parentCategory.Contain != "" {
		allowModels := yeeStrings.StringToIntArray(parentCategory.Contain, ",")
		if len(allowModels) != 0 {
			allow := false
			for _, v := range allowModels {
				if v == category.ModelId {
					allow = true
					break
				}
			}
			if !allow {
				return cateId, errors.New("此目录下不允许创建这个模型的栏目")
			}
		}
	}
	if category.ModelId == conf.DirModelId {
		containArr := yeeStrings.StringToIntArray(category.Contain, ",")
		if len(containArr) != 0 {
			category.Contain = yeeStrings.IntArrayToString(containArr, ",")
		}
	}
	if err := defaultDB.Create(&category).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return cateId, CreateError
	}
	cateId = category.Id
	return cateId, nil
}

// 编辑栏目
func (Category) Update(id int, name string) error {
	err := defaultDB.Model(&Category{Id: id}).Updates(Category{Name: name}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return DeleteError
	}
	return nil
}

// 排序栏目
func (Category) DoSort(id, sort int) error {
	err := defaultDB.Model(&Category{Id: id}).Updates(Category{Sort: sort}).Error
	if err != nil {
		easyweb.Logger.Error(err.Error())
		return SortError
	}
	return nil
}

// 删除栏目
func (Category) Delete(id int) error {
	if defaultDB.First(&Category{}, id).RecordNotFound() {
		return errors.New("栏目id错误")
	}
	subCategories := Category{}.SelectAllWithCache(id, conf.SuperAdminRoleId, false)
	if len(subCategories) > 0 {
		return errors.New("栏目下还有子栏目，请先删除子栏目")
	}
	if err := defaultDB.Delete(&Category{Id: id}).Error; err != nil {
		easyweb.Logger.Error(err.Error())
		return DeleteError
	}
	AdminPrivilege{}.Delete(id)
	SinglePage{}.Delete(id)
	return nil
}
