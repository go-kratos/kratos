package dao

import (
	"context"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/log"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/library/ecode"
)

// CreateTag 创建项目标签
func (d *Dao) CreateTag(c context.Context, tx *gorm.DB, pid int64, tagID string) error {

	tagMap := d.GetTagConfigInfo(c)
	covTagID, _ := strconv.ParseInt(tagID, 10, 64)

	if _, ok := tagMap[tagID]; !ok {
		//key不存在
		log.Error("标签id不存在对应标签名")
		tx.Rollback()
		return ecode.TicketAddTagFailed
	}

	// create
	if err := tx.Create(&model.ProjectTags{
		TagID:     covTagID,
		TagName:   tagMap[tagID],
		ProjectID: pid,
		Status:    1,
	}).Error; err != nil {
		log.Error("创建标签失败:%s", err)
		tx.Rollback()
		return err
	}

	return nil
}

// GetTagConfigInfo 获取标签配置信息返回map
func (d *Dao) GetTagConfigInfo(c context.Context) map[string]string {
	tagStr := conf.Conf.Tag.Tags
	allTags := strings.Split(tagStr, ",")

	tagMap := make(map[string]string)
	var tagInfo []string
	for _, v := range allTags {
		tagInfo = strings.Split(v, "=")
		tagMap[tagInfo[0]] = tagInfo[1]
	}
	return tagMap
}
