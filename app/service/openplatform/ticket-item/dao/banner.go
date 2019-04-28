package dao

import (
	"context"

	"encoding/json"
	"fmt"
	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
	"strconv"
	"strings"
	xtime "time"

	"github.com/jinzhu/gorm"
)

// DelBannerInfo 删除缓存所需banner信息
type DelBannerInfo struct {
	Order       int32
	Position    int32
	SubPosition int32
	DistrictIDs []string
}

// AddBanner 添加新投放
func (d *Dao) AddBanner(c context.Context, info *item.BannerEditRequest) (bannerID int64, verID uint64, err error) {
	// 开启事务
	tx := d.db.Begin()

	//创建banner信息
	banner := model.Banner{
		PubStart:   time.Time(info.PubStart),
		PubEnd:     time.Time(info.PubEnd),
		Name:       info.Name,
		Pic:        info.Pic,
		URL:        info.Url,
		From:       info.From,
		TargetUser: info.TargetUser,
	}

	if err = tx.Create(&banner).Error; err != nil {
		log.Error("创建banner失败:%s", err)
		tx.Rollback()
		return 0, 0, err
	}
	bannerID = banner.ID
	if bannerID == 0 {
		log.Error("bannerID为0")
		tx.Rollback()
		return 0, 0, err
	}

	var mainInfo []byte
	mainInfo, err = json.Marshal(info)
	if err != nil {
		log.Error("jsonMarshal版本详情失败:%s", err)
		tx.Rollback()
		return bannerID, 0, err
	}

	//生成新审核版本ver,verExt
	str := fmt.Sprintf("%d%d%.2d", info.Position, info.SubPosition, info.Order)
	forInt, _ := strconv.ParseInt(str, 10, 64)
	verInfo := model.Version{
		Type:       model.VerTypeBanner,
		Status:     info.OpType,
		ItemName:   info.Name,
		TargetItem: bannerID,
		AutoPub:    1, // 自动上架
		PubStart:   time.Time(info.PubStart),
		PubEnd:     time.Time(info.PubEnd),
		For:        forInt,
	}
	err = d.AddVersion(c, tx, &verInfo, &model.VersionExt{
		Type:     model.VerTypeBanner,
		MainInfo: string(mainInfo),
	})
	if err != nil {
		log.Error("创建banner版本失败: %s", err)
		return 0, verInfo.VerID, ecode.TicketAddVersionFailed
	}
	if verInfo.VerID == 0 {
		log.Error("创建后获取verID为0")
		return 0, verInfo.VerID, ecode.TicketAddVersionFailed
	}

	// 提交事务
	tx.Commit()
	return bannerID, verInfo.VerID, nil
}

// EditBanner 编辑投放
func (d *Dao) EditBanner(c context.Context, info *item.BannerEditRequest) (err error) {
	// 开启事务
	tx := d.db.Begin()

	var mainInfo []byte
	mainInfo, err = json.Marshal(info)
	if err != nil {
		log.Error("jsonMarshal版本详情失败:%s", err)
		tx.Rollback()
		return err
	}

	//更新审核版本ver,verExt
	str := fmt.Sprintf("%d%d%.2d", info.Position, info.SubPosition, info.Order)
	forInt, _ := strconv.ParseInt(str, 10, 64)
	if err = tx.Model(&model.Version{}).Where("ver_id = ?", info.VerId).Updates(
		map[string]interface{}{
			"item_name": info.Name,
			"status":    info.OpType,
			"pub_start": time.Time(info.PubStart),
			"pub_end":   time.Time(info.PubEnd),
			"for":       forInt,
		}).Error; err != nil {
		log.Error("更新banner版本失败:%s", err)
		tx.Rollback()
		return err
	}
	if err = tx.Model(&model.VersionExt{}).Where("ver_id = ?", info.VerId).Updates(
		map[string]interface{}{
			"main_info": mainInfo,
		}).Error; err != nil {
		log.Error("更新banner版本详情失败:%s", err)
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}

// PassOrPublishBanner 审核通过或上架banner版本
func (d *Dao) PassOrPublishBanner(c context.Context, verID uint64) (err error) {

	//获取版本信息
	verInfo, verExtInfo, getErr := d.GetVersion(c, verID, true)
	if getErr != nil {
		return getErr
	}

	//解析mainInfo
	var decodedMainInfo item.BannerEditRequest
	err = json.Unmarshal([]byte(verExtInfo.MainInfo), &decodedMainInfo)
	if err != nil {
		return err
	}
	bannerID := verInfo.TargetItem
	tx := d.db.Begin()

	//如果存在审核通过/进行中版本 删除版本
	var bannerVersions []model.Version
	if err = tx.Where("target_item = ? and status in (?) and ver_id != ? and deleted_at='0000-00-00 00:00:00'",
		bannerID, []int32{model.VerStatusReadyForSale, model.VerStatusOnShelf}, verID).Find(&bannerVersions).Error; err != nil {
		log.Error("获取banner所有审核通过/已上架版本失败:%s", err)
		tx.Rollback()
		return err
	}

	var delCacheInfos []DelBannerInfo
	var verIDs []uint64
	var delPosition int32
	var delSubPosition int32
	var delOrder int32
	var delDistrictIDs []string
	for _, v := range bannerVersions {
		//获取该版本详情 删除对应bannerDistrict信息
		if delPosition, delSubPosition, delOrder, delDistrictIDs, err = d.DelBannerDistrictByVerID(c, tx, v.VerID, bannerID); err != nil {
			return
		}
		delCacheInfos = append(delCacheInfos, DelBannerInfo{
			Order:       delOrder,
			Position:    delPosition,
			SubPosition: delSubPosition,
			DistrictIDs: delDistrictIDs,
		})
		verIDs = append(verIDs, v.VerID)
	}
	if verIDs != nil {
		if err = tx.Exec("UPDATE version SET deleted_at=? WHERE ver_id in (?)", xtime.Now().Format("2006-01-02 15:04:05"), verIDs).Error; err != nil {
			log.Error("删除banner版本失败:%s", err)
			tx.Rollback()
			return err
		}
		if err = tx.Exec("UPDATE version_ext SET deleted_at=? WHERE ver_id in (?)", xtime.Now().Format("2006-01-02 15:04:05"), verIDs).Error; err != nil {
			log.Error("删除banner ext版本失败:%s", err)
			tx.Rollback()
			return err
		}
	}

	//如果达到投放开始时间 还要更新bannerDistrict并直接上架
	var newVerStatus int32
	var newBannerStatus int32
	var districtIDs []string
	if decodedMainInfo.PubStart <= xtime.Now().Unix() {
		districtIDs = strings.Split(decodedMainInfo.Location, ",")
		for _, districtID := range districtIDs {
			newDistrictID, _ := strconv.ParseInt(districtID, 10, 64)
			if err = d.CreateOrUpdateBannerDistrict(c, tx, model.BannerDistrict{
				BannerID:    bannerID,
				Position:    decodedMainInfo.Position,
				SubPosition: decodedMainInfo.SubPosition,
				Order:       decodedMainInfo.Order,
				DistrictID:  newDistrictID,
			}); err != nil {
				return err
			}
		}
		newVerStatus = model.VerStatusOnShelf
		newBannerStatus = 1
	} else {
		newVerStatus = model.VerStatusReadyForSale
		newBannerStatus = 0
	}

	//更新banner表
	if err = tx.Model(&model.Banner{}).Where("id = ?", bannerID).Updates(
		map[string]interface{}{
			"pub_start":   time.Time(decodedMainInfo.PubStart),
			"pub_end":     time.Time(decodedMainInfo.PubEnd),
			"name":        decodedMainInfo.Name,
			"pic":         decodedMainInfo.Pic,
			"url":         decodedMainInfo.Url,
			"from":        decodedMainInfo.From,
			"status":      newBannerStatus,
			"target_user": decodedMainInfo.TargetUser,
		}).Error; err != nil {
		log.Error("更新banner失败:%s", err)
		tx.Rollback()
		return err
	}

	//更新version状态为审核通过(待上架)或已上架
	if err = tx.Model(&model.Version{}).Where("ver_id = ?", verID).Updates(
		map[string]interface{}{
			"status": newVerStatus,
		}).Error; err != nil {
		log.Error("更新banner审核版本失败:%s", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	//如果新状态是上架状态 删除相关缓存
	if newVerStatus == model.VerStatusOnShelf {
		_, err = d.DelBannerCache(c, decodedMainInfo.Position, decodedMainInfo.SubPosition, decodedMainInfo.Order, districtIDs, bannerID)
		if err != nil {
			return
		}
		for _, v := range delCacheInfos {
			_, err = d.DelBannerCache(c, v.Position, v.SubPosition, v.Order, v.DistrictIDs, 0)
			if err != nil {
				return
			}
		}

	}

	return
}

// DelBannerCache 删除banner相关缓存
func (d *Dao) DelBannerCache(c context.Context, position int32, subPosition int32, order int32, districtIDs []string, bannerID int64) (res bool, err error) {
	var (
		keys []interface{}
	)

	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()

	for _, districtID := range districtIDs {
		// DEL bannerList
		keys = append(keys, keyBannerList(order, districtID, position, subPosition))
	}
	//删除bannerInfo缓存
	if bannerID != 0 {
		keys = append(keys, keyBannerInfo(bannerID))
	}

	log.Info("DEL %v", keys)
	//DEL bannerList
	if err = conn.Send("DEL", keys...); err != nil {
		log.Error("DEL %v, error(%v)", keys, err)
	}
	if err != nil {
		return false, err
	}
	return true, err
}

//DelBannerDistrictByVerID 删除版本详情内对应的关系表信息
func (d *Dao) DelBannerDistrictByVerID(c context.Context, tx *gorm.DB, verID uint64, bannerID int64) (position int32, subPosition int32, order int32, districtIDs []string, err error) {
	var verExtInfo model.VersionExt
	var delDecodedInfo item.BannerEditRequest
	//获取该版本详情 删除对应bannerDistrict信息
	if err = tx.Where("ver_id = ?", verID).First(&verExtInfo).Error; err != nil {
		log.Error("获取需要删除版本详情失败:%s", err)
		tx.Rollback()
		return 0, 0, 0, nil, err
	}
	//解析mainInfo
	err = json.Unmarshal([]byte(verExtInfo.MainInfo), &delDecodedInfo)
	if err != nil {
		tx.Rollback()
		return 0, 0, 0, nil, err
	}

	districtIDs = strings.Split(delDecodedInfo.Location, ",")
	for _, districtID := range districtIDs {
		if err = tx.Exec("UPDATE banner_district SET is_deleted=1 WHERE banner_id = ? and district_id = ? and position = ? and sub_position = ? and `order` = ?",
			bannerID, districtID, delDecodedInfo.Position, delDecodedInfo.SubPosition, delDecodedInfo.Order).Error; err != nil {
			log.Error("删除bannerDistrict失败:%s", err)
			tx.Rollback()
			return 0, 0, 0, nil, err
		}
	}
	return delDecodedInfo.Position, delDecodedInfo.SubPosition, delDecodedInfo.Order, districtIDs, err
}

// CreateOrUpdateBannerDistrict 更新或新建bannerDistrict关系记录
func (d *Dao) CreateOrUpdateBannerDistrict(c context.Context, tx *gorm.DB, info model.BannerDistrict) (err error) {
	var bannerDist model.BannerDistrict
	log.Info("bannerDist:%v", info)
	err = tx.Where("district_id = ? and position = ? and sub_position = ? and `order` = ?",
		info.DistrictID, info.Position, info.SubPosition, info.Order).First(&bannerDist).Error
	//除去没查找到记录的报错 其他直接抛错
	if err != nil && err != ecode.NothingFound {
		log.Error("获取banner dist信息失败:%s", err)
		tx.Rollback()
		return
	}

	if bannerDist.ID != 0 {
		//update
		log.Info("update bannerDistrict")
		if err = tx.Exec("UPDATE banner_district SET banner_id=?,is_deleted=0 WHERE district_id = ? and position = ? and sub_position = ? and `order` = ?",
			info.BannerID, info.DistrictID, info.Position, info.SubPosition, info.Order).Error; err != nil {
			log.Error("更新bannerDistrict失败:%s", err)
			tx.Rollback()
			return err
		}
	} else {
		//create
		log.Info("insert bannerDistrict")
		if err = tx.Create(&model.BannerDistrict{
			BannerID:    info.BannerID,
			Position:    info.Position,
			SubPosition: info.SubPosition,
			Order:       info.Order,
			DistrictID:  info.DistrictID,
		}).Error; err != nil {
			log.Error("创建bannerDistrict失败:%s", err)
			tx.Rollback()
			return err
		}
	}
	return nil
}

// DeleteBanner 删除banner
func (d *Dao) DeleteBanner(c context.Context, verID uint64) (err error) {
	tx := d.db.Begin()
	if err = tx.Exec("UPDATE version SET deleted_at=? WHERE ver_id = ?", xtime.Now().Format("2006-01-02 15:04:05"), verID).Error; err != nil {
		tx.Rollback()
		log.Error("删除banner版本失败:%s", err)
		return err
	}
	if err = tx.Exec("UPDATE version_ext SET deleted_at=? WHERE ver_id=?", xtime.Now().Format("2006-01-02 15:04:05"), verID).Error; err != nil {
		tx.Rollback()
		log.Error("删除banner ext版本失败:%s", err)
		return err
	}
	tx.Commit()
	return
}

// UnpublishBannerManual 手动取消激活banner
func (d *Dao) UnpublishBannerManual(c context.Context, verID uint64) (err error) {
	//取消激活banner
	//获取版本信息
	tx := d.db.Begin()
	verInfo, _, getErr := d.GetVersion(c, verID, true)
	if getErr != nil {
		return getErr
	}
	bannerID := verInfo.TargetItem
	if err = tx.Exec("UPDATE banner SET status = 0 WHERE id = ?", bannerID).Error; err != nil {
		tx.Rollback()
		log.Error("更新banner未激活状态失败:%s", err)
		return err
	}

	//更新版本为草稿状态
	if err = tx.Exec("UPDATE version SET status = ? WHERE ver_id = ?", model.VerStatusNotReviewed, verID).Error; err != nil {
		tx.Rollback()
		log.Error("更新banner版本为草稿状态失败:%s", err)
		return err
	}

	//删除此版本bannerDistrict信息
	var delPosition int32
	var delSubPosition int32
	var delOrder int32
	var delDistrictIDs []string
	if delPosition, delSubPosition, delOrder, delDistrictIDs, err = d.DelBannerDistrictByVerID(c, tx, verID, bannerID); err != nil {
		return
	}

	//如果存在除这个版本以外同个投放id的版本 删除版本
	var bannerVersions []model.Version
	if err = tx.Where("target_item = ? and ver_id != ? and deleted_at='0000-00-00 00:00:00'",
		bannerID, verID).Find(&bannerVersions).Error; err != nil {
		log.Error("获取banner所有其他版本失败:%s", err)
		tx.Rollback()
		return err
	}

	var verIDs []uint64
	for _, v := range bannerVersions {
		//获取该版本详情 删除对应bannerDistrict信息
		_, _, _, _, err = d.DelBannerDistrictByVerID(c, tx, v.VerID, bannerID)
		if err != nil {
			return
		}
		verIDs = append(verIDs, v.VerID)
	}
	if verIDs != nil {
		if err = tx.Exec("UPDATE version SET deleted_at=? WHERE ver_id in (?)", xtime.Now().Format("2006-01-02 15:04:05"), verIDs).Error; err != nil {
			log.Error("删除banner版本失败:%s", err)
			tx.Rollback()
			return err
		}
		if err = tx.Exec("UPDATE version_ext SET deleted_at=? WHERE ver_id in (?)", xtime.Now().Format("2006-01-02 15:04:05"), verIDs).Error; err != nil {
			log.Error("删除banner ext版本失败:%s", err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	//删除下架相关缓存
	_, err = d.DelBannerCache(c, delPosition, delSubPosition, delOrder, delDistrictIDs, bannerID)

	return
}

// UnpublishBannerForced 已过期自动下架操作
func (d *Dao) UnpublishBannerForced(c context.Context, verID uint64) (err error) {
	//获取版本信息
	tx := d.db.Begin()
	verInfo, _, getErr := d.GetVersion(c, verID, true)
	if getErr != nil {
		return getErr
	}
	//取消激活banner
	bannerID := verInfo.TargetItem
	if err = tx.Exec("UPDATE banner SET status = 0 WHERE id = ?", bannerID).Error; err != nil {
		tx.Rollback()
		log.Error("更新banner未激活状态失败:%s", err)
		return err
	}
	//删除bannerDistrict信息
	var delPosition int32
	var delSubPosition int32
	var delOrder int32
	var delDistrictIDs []string
	if delPosition, delSubPosition, delOrder, delDistrictIDs, err = d.DelBannerDistrictByVerID(c, tx, verID, bannerID); err != nil {
		return
	}

	//更新版本为强制下架状态
	if err = tx.Exec("UPDATE version SET status = ? WHERE ver_id = ?", model.VerStatusOffShelfForced, verID).Error; err != nil {
		tx.Rollback()
		log.Error("更新banner版本为强制下架状态失败:%s", err)
		return err
	}
	tx.Commit()
	//删除下架相关缓存
	_, err = d.DelBannerCache(c, delPosition, delSubPosition, delOrder, delDistrictIDs, bannerID)

	return
}

// CgBannerStatus 更改banner版本状态
func (d *Dao) CgBannerStatus(c context.Context, info *item.VersionStatusRequest) (err error) {

	switch info.OpType {
	case 0:
		//手动取消激活操作
		err = d.UnpublishBannerManual(c, info.VerId)
	case 1:
		//提交审核操作
		err = d.db.Exec("UPDATE version SET status = ? WHERE ver_id = ?", model.VerStatusReadyForReview, info.VerId).Error
		if err != nil {
			return err
		}
		//提交审核时 记入versionLog
		err = d.AddVersionLog(c, &model.VersionLog{
			VerID: info.VerId,
			Type:  2, //用户操作记录
			Log:   "提交审核",
			Uname: info.Uname,
		})
	case 2:
		//删除操作
		err = d.DeleteBanner(c, info.VerId)
	case 3:
		//上架操作
		err = d.PassOrPublishBanner(c, info.VerId)
	case 4:
		//已过期自动下架操作
		err = d.UnpublishBannerForced(c, info.VerId)
	default:
		return ecode.NothingFound
	}

	return
}
