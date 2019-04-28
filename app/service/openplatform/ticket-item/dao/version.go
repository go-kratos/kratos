package dao

import (
	"context"

	"fmt"
	"github.com/jinzhu/gorm"
	"go-common/app/common/openplatform/random"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddVersion 添加新版本
func (d *Dao) AddVersion(c context.Context, requiredTx *gorm.DB, verInfo *model.Version, verExtInfo *model.VersionExt) (err error) {
	// 此处ver_id调用订单号生成器
	verID := uint64(random.Uniqid(19))

	// 版本号赋值与信息
	verInfo.VerID = verID
	verExtInfo.VerID = verID

	var tx *gorm.DB
	if requiredTx == nil {
		// 开启事务
		tx = d.db.Begin()
	} else {
		tx = requiredTx
	}

	// 插入新数据
	if verError := tx.Create(&verInfo).Error; verError != nil {
		log.Error("version insertion failed:%s", verError)
		tx.Rollback()
		return ecode.TicketAddVersionFailed
	}

	if extError := tx.Create(&verExtInfo).Error; extError != nil {
		log.Error("version ext insertion failed:%s", extError)
		tx.Rollback()
		return ecode.TicketAddVerExtFailed
	}

	if requiredTx == nil {
		tx.Commit()
	}
	return nil
}

// UpdateVersion 编辑版本信息
func (d *Dao) UpdateVersion(c context.Context, verInfo *model.Version) (bool, error) {

	// update guest with new info (using map can update the column with empty string)
	updateErr := d.db.Model(&model.Version{}).Where("ver_id = ?", verInfo.VerID).Updates(
		map[string]interface{}{
			"type":        verInfo.Type,
			"status":      verInfo.Status,
			"item_name":   verInfo.ItemName,
			"ver":         verInfo.Ver,
			"target_item": verInfo.TargetItem,
			"auto_pub":    verInfo.AutoPub,
			"parent_id":   verInfo.ParentID,
		}).Error
	if updateErr != nil {
		log.Error("VERSION UPDATE FAILED:%s", updateErr)
		return false, ecode.NotModified
	}

	return true, nil
}

// GetVersion 获取版本信息外加详情
func (d *Dao) GetVersion(c context.Context, verID uint64, needExt bool) (*model.Version, *model.VersionExt, error) {

	var verInfo model.Version
	var verExtInfo model.VersionExt

	if dbErr := d.db.Where("ver_id = ?", verID).First(&verInfo).Error; dbErr != nil {
		log.Error("verinfo:(%v) not found with err:%s", verID, dbErr)
		return nil, nil, ecode.NothingFound
	}

	if needExt {
		if dbErr := d.db.Where("ver_id = ?", verID).First(&verExtInfo).Error; dbErr != nil {
			log.Error("ver_ext_info:(%v) not found with err:%s", verID, dbErr)
			return nil, nil, ecode.NothingFound
		}
	}

	return &verInfo, &verExtInfo, nil
}

// RejectVersion 驳回版本
func (d *Dao) RejectVersion(c context.Context, verID uint64, verType int32) (bool, error) {

	var newStatus int32
	switch verType {
	case model.VerTypeBanner:
		newStatus = model.VerStatusNotReviewed
	default:
		newStatus = model.VerStatusRejected
	}
	updateErr := d.db.Where("ver_id = ? and type = ?", verID, verType).Model(&model.Version{}).Update("status", newStatus).Error
	if updateErr != nil {
		log.Error("更新版本状态失败:%s", updateErr)
		return false, ecode.NotModified
	}

	return true, nil
}

// AddVersionLog 新建版本审核记录
func (d *Dao) AddVersionLog(c context.Context, info *model.VersionLog) error {
	if insertErr := d.db.Create(&info).Error; insertErr != nil {
		log.Error("新建版本审核记录失败:%s", insertErr)
		return ecode.NotModified
	}
	return nil
}

// VersionSearch 项目版本查询
func (d *Dao) VersionSearch(c context.Context, in *model.VersionSearchParam) (versions *model.VersionSearchList, err error) {
	r := d.es.NewRequest("ticket_version").Index("ticket_version")
	if in.TargetItem > 0 {
		r.WhereEq("target_item", in.TargetItem)
	} else if in.ItemName != "" {
		r.WhereLike([]string{"item_name"}, []string{in.ItemName}, false, elastic.LikeLevelLow)
	}
	if in.Type > 0 {
		r.WhereEq("type", in.Type)
	}
	if length := len(in.Status); length == 1 {
		r.WhereEq("status", in.Status[0])
	} else if length > 1 {
		r.WhereIn("status", in.Status)
	}
	r.Order("ctime", elastic.OrderDesc).Ps(in.Ps).Pn(in.Pn)

	log.Info(fmt.Sprintf("%s/x/admin/search/query?%s", d.c.URL.ElasticHost, r.Params()))

	versions = new(model.VersionSearchList)
	err = r.Scan(c, versions)
	if err != nil {
		log.Error("VersionSearch(%v) r.Query(%s) error(%s)", in, r.Params(), err)
		return
	}
	return
}
