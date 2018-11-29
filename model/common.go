/**
 * Created by WillkYang on 2018/10/9.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */
package model

func DoSort(modelId, id, value int) error {
	m := Model{Id: modelId}
	if err := m.SelectOneModel(&m); err != nil {
		return err
	}
	return defaultDB.Table(m.DbTableName).Where("id = ?", id).Update("sort", value).Error
}
