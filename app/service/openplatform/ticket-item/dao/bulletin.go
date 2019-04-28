package dao

import (
	"context"

	"encoding/json"
	"go-common/app/common/openplatform/random"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
)

// BulletinMainInfo 公告版本内容
type BulletinMainInfo struct {
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
	Content      string `json:"content"`
	ProjectName  string `json:"project_name"`
	Pid          int64  `json:"pid"`
}

// GetBulletins 获取项目下所有公告
func (d *Dao) GetBulletins(c context.Context, pid int64) (res []*item.BulletinInfo, err error) {
	// 获取项目下所有bulletin基础信息
	var projectBulletins []*model.Bulletin
	if dbErr := d.db.Order("ctime desc").Where("project_id = ? and status = 1", pid).Find(&projectBulletins).Error; dbErr != nil {
		log.Error("获取项目公告基础信息失败: %s", dbErr)
		return nil, ecode.NothingFound
	}

	var bulIDs []int64
	for _, value := range projectBulletins {
		bulIDs = append(bulIDs, value.ID)
	}

	// 获取公告详情并用id当键值组成map
	var bulletinDetail []*model.BulletinExtra
	if dbErr := d.db.Where("id in (?)", bulIDs).Find(&bulletinDetail).Error; dbErr != nil {
		log.Error("获取项目公告详情失败: %s", dbErr)
		return nil, ecode.NothingFound
	}
	bulletinDetailMap := make(map[int64]*model.BulletinExtra)
	for _, value := range bulletinDetail {
		bulletinDetailMap[value.ID] = value
	}

	// 添加详情信息到公告信息中
	for _, value := range projectBulletins {
		tmpBul := &item.BulletinInfo{
			ID:      value.ID,
			Title:   value.Title,
			Content: value.Content,
			Ctime:   value.Ctime.Time().Format("2006-01-02 15:04:05"),
			Mtime:   value.Mtime.Time().Format("2006-01-02 15:04:05"),
			VerID:   value.VerID,
		}
		detailInfo, ok := bulletinDetailMap[value.ID]
		if ok {
			tmpBul.Detail = detailInfo.Detail
		} else {
			tmpBul.Detail = ""
		}
		res = append(res, tmpBul)
	}

	return
}

// AddBulletin 添加公告
func (d *Dao) AddBulletin(c context.Context, info *item.BulletinInfoRequest) (bool, error) {
	pid := info.ParentID
	mainInfo, jsonErr := d.GenBulMainInfo(pid, info.Title, info.Content, info.Detail)

	if jsonErr != nil {
		log.Error("获取整合maininfo失败: %s", jsonErr)
		return false, ecode.NothingFound
	}

	// add version and version ext
	verErr := d.AddVersion(c, nil, &model.Version{
		Type:       model.VerTypeBulletin,
		Status:     1, // 审核中
		ItemName:   info.Title,
		ParentID:   pid,
		TargetItem: info.TargetItem,
		AutoPub:    1, // 自动上架
		PubStart:   model.TimeNull,
		PubEnd:     model.TimeNull,
	}, &model.VersionExt{
		Type:     model.VerTypeBulletin,
		MainInfo: string(mainInfo),
	})
	if verErr != nil {
		log.Error("创建公告版本失败: %s", verErr)
		return false, ecode.NothingFound
	}

	return true, nil
}

// PassBulletin 审核通过公告
func (d *Dao) PassBulletin(c context.Context, verID uint64) (bool, error) {
	verInfo, verExtInfo, err := d.GetVersion(c, verID, true)
	if err != nil {
		return false, ecode.NothingFound
	}
	targetItem := verInfo.TargetItem

	var decodedMainInfo BulletinMainInfo
	err = json.Unmarshal([]byte(verExtInfo.MainInfo), &decodedMainInfo)
	if err != nil {
		return false, err
	}

	var finalTargetItem int64

	// 开启事务
	tx := d.db.Begin()
	if targetItem == 0 {
		// 没新建过bulletin信息
		bulletinID := random.Uniqid(19)
		bulData := &model.Bulletin{
			Status:     1,
			Title:      decodedMainInfo.Name,
			Content:    decodedMainInfo.Introduction,
			ProjectID:  decodedMainInfo.Pid,
			VerID:      verID,
			BulletinID: bulletinID,
		}
		insertErr := tx.Save(&bulData).Error
		if insertErr != nil {
			tx.Rollback()
			log.Error("新建bulletin失败: %s", insertErr)
			return false, ecode.NotModified
		}

		// 获取新建的bulletin自增id
		bulPrimID := bulData.ID
		if bulPrimID == 0 {
			tx.Rollback()
			log.Error("获取新建bulletin自增id失败")
			return false, ecode.NothingFound
		}

		// 新建bulletin_extra
		insertExtErr := tx.Create(&model.BulletinExtra{
			ID:         bulPrimID,
			Detail:     decodedMainInfo.Content,
			BulletinID: bulletinID,
		}).Error
		if insertExtErr != nil {
			tx.Rollback()
			log.Error("新建bulletin_extra失败：%s", insertExtErr)
			return false, ecode.NotModified
		}

		finalTargetItem = bulPrimID

	} else {
		// 已建过的直接更新bulletin
		updateErr := tx.Where("id = ?", targetItem).Model(&model.Bulletin{}).Updates(
			map[string]interface{}{
				"status":  1,
				"title":   decodedMainInfo.Name,
				"content": decodedMainInfo.Introduction,
				"ver_id":  verID,
			}).Error
		if updateErr != nil {
			tx.Rollback()
			log.Error("UPDATE BULLETIN FAILED")
			return false, ecode.NotModified
		}
		// 更新bulletin_extra
		updateExtErr := tx.Where("id = ?", targetItem).Model(&model.BulletinExtra{}).Updates(
			map[string]interface{}{
				"detail": decodedMainInfo.Content,
			}).Error
		if updateExtErr != nil {
			tx.Rollback()
			log.Error("UPDATE BULLETIN_EXTRA FAILED")
			return false, ecode.NotModified
		}
		finalTargetItem = targetItem

	}
	// 更新版本为已审核状态
	updateVerErr := tx.Where("ver_id = ?", verID).Model(&model.Version{}).Updates(
		map[string]interface{}{
			"status":      4, // 已上架状态
			"ver":         "1.0",
			"target_item": finalTargetItem,
		}).Error
	if updateVerErr != nil {
		tx.Rollback()
		log.Error("UPDATE BULLETIN VERSION FAILED")
		return false, ecode.NotModified
	}

	tx.Commit()
	return true, nil
}

// UpdateBulletin 编辑公告版本
func (d *Dao) UpdateBulletin(c context.Context, info *item.BulletinInfoRequest) (bool, error) {

	// 获取JSONEncode好的maininfo
	pid := info.ParentID
	mainInfo, jsonErr := d.GenBulMainInfo(pid, info.Title, info.Content, info.Detail)

	if jsonErr != nil || mainInfo == "" {
		log.Error("获取maininfo失败: %s", jsonErr)
		return false, ecode.NothingFound
	}

	// 开启事务
	tx := d.db.Begin()

	// 编辑version_ext
	updateExtErr := tx.Model(&model.VersionExt{}).Where("ver_id = ? and type = ?", info.VerID, model.VerTypeBulletin).Update("main_info", mainInfo).Error
	if updateExtErr != nil {
		tx.Rollback()
		log.Error("更新versionext失败:%s", updateExtErr)
		return false, ecode.NotModified
	}

	// 编辑version
	updateErr := d.db.Model(&model.Version{}).Where("ver_id = ? and type = ?", info.VerID, model.VerTypeBulletin).Updates(
		map[string]interface{}{
			"status":    1, // 审核中
			"item_name": info.Title,
		}).Error
	if updateErr != nil {
		tx.Rollback()
		log.Error("VERSION UPDATE FAILED:%s", updateErr)
		return false, ecode.NotModified
	}

	tx.Commit()

	return true, nil
}

// GenBulMainInfo 整合公告的详情maininfo字段
func (d *Dao) GenBulMainInfo(pid int64, name string, introduction string, content string) (string, error) {

	var projectInfo []model.Item
	if dbErr := d.db.Select("name").Where("id = ?", pid).Find(&projectInfo).Error; dbErr != nil {
		log.Error("获取项目信息失败:%s", dbErr)
		return "", ecode.NothingFound
	}
	extInfo := BulletinMainInfo{
		Name:         name,
		Introduction: introduction,
		Content:      content,
		ProjectName:  projectInfo[0].Name,
		Pid:          pid,
	}
	mainInfo, jsonErr := json.Marshal(extInfo)

	if jsonErr != nil {
		log.Error("JSONEncode失败: %s", jsonErr)
		return "", ecode.NothingFound
	}
	return string(mainInfo), nil
}

// UnpublishBulletin 下架版本
func (d *Dao) UnpublishBulletin(c context.Context, verID uint64, status int8) (bool, error) {

	// 获取版本信息
	verInfo, _, err := d.GetVersion(c, verID, false)
	if err != nil {
		return false, ecode.NothingFound
	}

	bulletinID := verInfo.TargetItem

	// 开启事务
	tx := d.db.Begin()
	// 取消激活公告
	updateErr := tx.Model(&model.Bulletin{}).Where("id = ? and ver_id = ?", bulletinID, verID).Updates(
		map[string]interface{}{
			"status": 0,
		}).Error
	if updateErr != nil {
		tx.Rollback()
		log.Error("取消激活公告失败:%s", updateErr)
		return false, ecode.NotModified
	}
	// 更新版本为已下架状态
	updateVerErr := tx.Model(&model.Version{}).Where("ver_id = ? and type = ?", verID, model.VerTypeBulletin).Updates(
		map[string]interface{}{
			"status": status,
		}).Error
	if updateVerErr != nil {
		tx.Rollback()
		log.Error("版本更新状态失败:%s", updateVerErr)
	}

	tx.Commit()
	return true, nil

}
