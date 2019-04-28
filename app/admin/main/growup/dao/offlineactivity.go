package dao

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/growup/dao/shell"
	"go-common/app/admin/main/growup/model/offlineactivity"
	"go-common/app/admin/main/up/util"
	"go-common/library/log"
	"strconv"
	"strings"
)

func generateDelimiter(delimiter string) func(r rune) bool {
	return func(r rune) bool {
		for _, v := range delimiter {
			if v == r {
				return true
			}
		}
		return false
	}
}

const (
	trimSet = "\r\n "
)

func strlistToInt64List(list []string, trimSet string) (result []int64) {
	for _, v := range list {
		mid, e := strconv.ParseInt(strings.Trim(v, trimSet), 10, 64)
		if e != nil {
			continue
		}
		result = append(result, mid)
	}
	return
}

func trimString(strlist []string, trimset string) (result []string) {
	for _, v := range strlist {
		var trimed = strings.Trim(v, trimSet)
		if trimed == "" {
			continue
		}
		result = append(result, trimed)
	}
	return
}

//ParseMidsFromString parse []int64 from string
func ParseMidsFromString(str string) (result []int64) {
	var midstrlist = trimString(strings.FieldsFunc(str, generateDelimiter(",\r\n")), trimSet)
	result = util.Unique(strlistToInt64List(midstrlist, trimSet))
	return
}

//OfflineActivityAddActivity add acitvity
func (d *Dao) OfflineActivityAddActivity(ctx context.Context, arg *offlineactivity.AddActivityArg) (err error) {
	var activityInfo = offlineactivity.OfflineActivityInfo{
		Title:     arg.Title,
		Link:      arg.Link,
		BonusType: arg.BonusType,
		Memo:      arg.Memo,
		State:     int8(offlineactivity.ActivityStateInit),
		Creator:   arg.Creator,
	}
	var bonusList = arg.BonusList
	var hasBonus = false
	for _, bonus := range bonusList {
		if bonus.TotalMoney <= 0 {
			err = fmt.Errorf("bonus money < 0, money=%f", bonus.TotalMoney)
			log.Error(err.Error())
			return
		}
		bonus.MidList = ParseMidsFromString(bonus.Mids)
		bonus.MemberCount = int64(len(bonus.MidList))
		if bonus.MemberCount == 0 {
			log.Warn("no mid for this bonus, bonus：money=%.2f", bonus.TotalMoney)
			continue
		}
		hasBonus = true
		log.Info("bonus info:money=%.2f, membercount=%d", bonus.TotalMoney/1000, bonus.MemberCount)
	}
	// 如果一个bonus都没有，则什么也不做
	if !hasBonus {
		log.Error("no bonus")
		err = fmt.Errorf("no bonus info")
		return
	}
	if err = d.db.Create(&activityInfo).Error; err != nil {
		log.Error("err insert offline activity, err=%s", err)
		return
	}

	var tx = d.db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
			activityInfo.State = int8(offlineactivity.ActivityStateCreateFail)
			if err = d.db.Select("state").Save(&activityInfo).Error; err != nil {
				log.Error("err update offline activity, err=%s", err)
				return
			}
		}
	}()
	// create

	for _, bonus := range bonusList {
		if bonus.MemberCount <= 0 {
			continue
		}
		// 插入OfflineActivityBonus
		//if activityInfo.BonusType == offlineactivity.BonusTypeMoney {
		bonusMoney := offlineactivity.GetMoneyForDb(bonus.TotalMoney)
		//}
		var bonusInfo = offlineactivity.OfflineActivityBonus{
			TotalMoney:  bonusMoney,
			MemberCount: uint32(bonus.MemberCount),
			State:       int8(offlineactivity.BonusStateInit),
			ActivityID:  activityInfo.ID,
		}
		if err = tx.Create(&bonusInfo).Error; err != nil {
			log.Error("err insert offline activity bonus, err=%s", err)
			return
		}

		// 插入OfflineActivityResult
		var insertSchema []string
		var vals []interface{}
		var dbmoney = bonusInfo.TotalMoney / int64(bonusInfo.MemberCount)
		for _, mid := range bonus.MidList {
			insertSchema = append(insertSchema, "(?,?,?,?,?)")
			vals = append(vals, activityInfo.ID, bonusInfo.ID, activityInfo.BonusType, mid, dbmoney)
		}
		var insertPre = "insert into offline_activity_result (activity_id, bonus_id, bonus_type, mid, bonus_money) values "
		var insertSQL = insertPre + strings.Join(insertSchema, ",")
		if err = tx.Table(offlineactivity.TableOfflineActivityResult).Exec(insertSQL, vals...).Error; err != nil {
			log.Error("err insert offline activity result, err=%s", err)
			return
		}
	}

	activityInfo.State = int8(offlineactivity.ActivityStateSending)
	if arg.BonusType == int8(offlineactivity.BonusTypeThing) {
		activityInfo.State = int8(offlineactivity.ActivityStateSucess)
	}
	if err = tx.Select("state").Save(&activityInfo).Error; err != nil {
		log.Error("err update offline activity, err=%s", err)
		return
	}
	return tx.Commit().Error
}

//ShellCallbackUpdate shell callback
func (d *Dao) ShellCallbackUpdate(ctx context.Context, result *shell.OrderCallbackJSON, msgid string) (orderInfo offlineactivity.OfflineActivityShellOrder, err error) {

	// 找到order info, 写入结果
	var stateResult = offlineactivity.ActivityStateWaitResult
	if result.IsSuccess() {
		stateResult = offlineactivity.ActivityStateSucess
	} else if result.IsFail() {
		stateResult = offlineactivity.ActivityStateFail
	}

	if err = d.db.Where("order_id=?", result.ThirdOrderNo).Find(&orderInfo).Error; err != nil {
		log.Error("order id=%s, find order id fail, err=%s", result.ThirdOrderNo, err)
		return
	}

	var tx = d.db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Table(offlineactivity.TableOfflineActivityShellOrder).
		Where("id=?", orderInfo.ID).
		Update("order_status", result.Status).
		Error; err != nil {
		log.Error("order id=%s, update order fail, err=%s", result.ThirdOrderNo, err)
		return
	}

	if err = tx.Table(offlineactivity.TableOfflineActivityResult).
		Where("id=?", orderInfo.ResultID).
		Update("state", stateResult).
		Error; err != nil {
		log.Error("order id=%s, update result fail, err=%s", result.ThirdOrderNo, err)
		return
	}
	log.Info("order id=%s, shell callback, success, mid=%s, status=%s", result.ThirdOrderNo, result.Mid, result.Status)

	err = tx.Commit().Error
	return
}

//OfflineActivityGetUpBonusResult Up主所有Activity一起统计
func (d *Dao) OfflineActivityGetUpBonusResult(ctx context.Context, needCount bool, limit, offset int, query string, args ...interface{}) (upResult []*offlineactivity.OfflineActivityResult, totalCount int, err error) {
	// 查询所有的 result for mid
	// 区分已结算、未结算
	// 已结算= 成功， 未结算= 初始、发送、等待
	if needCount {
		if err = d.db.Table(offlineactivity.TableOfflineActivityResult).Where(query, args...).Count(&totalCount).Error; err != nil {
			log.Error("fail to get from db, err=%s", err)
			return
		}
	}

	if err = d.db.Where(query, args...).Limit(limit).Offset(offset).Find(&upResult).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	return
}

//OfflineActivityGetUpBonusResultSelect select by selection
func (d *Dao) OfflineActivityGetUpBonusResultSelect(ctx context.Context, selectQuery string, query string, args ...interface{}) (upResult []*offlineactivity.OfflineActivityResult, err error) {
	if err = d.db.Where(query, args...).Select(selectQuery).Find(&upResult).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	return
}

//OfflineActivityGetUpBonusByActivityResult 根据Activity来分别统计
func (d *Dao) OfflineActivityGetUpBonusByActivityResult(ctx context.Context, limit, offset int, mid int64) (upResult []*offlineactivity.OfflineActivityResult, totalCount int, err error) {
	// 查询所有的 result for mid
	// 区分已结算、未结算
	// 已结算= 成功， 未结算= 初始、发送、等待
	var actividyIds []struct {
		ActivityID int64 `gorm:"column:activity_id"`
	}

	if err = d.db.Table(offlineactivity.TableOfflineActivityResult).Select("count(distinct(activity_id))").Where("mid=?", mid).Count(&totalCount).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	if err = d.db.Table(offlineactivity.TableOfflineActivityResult).Select("distinct(activity_id)").Where("mid=?", mid).Offset(offset).Limit(limit).Find(&actividyIds).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	//totalCount = len(actividyIds)
	var ids []int64
	for i := 0; i < len(actividyIds); i++ {
		ids = append(ids, actividyIds[i].ActivityID)
	}

	if err = d.db.Where("mid=? and activity_id in(?)", mid, ids).Find(&upResult).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	return
}

//OfflineActivityGetDB get gorm db
func (d *Dao) OfflineActivityGetDB() *gorm.DB {
	return d.db
}

//UpdateActivityState 更新Activity的状态，根据每一个付款记录的状态
func (d *Dao) UpdateActivityState(ctx context.Context, activityID int64) (affectedRow int64, err error) {
	// 找到所有result中，state最小的一个
	// 更新到activity的状态
	var record offlineactivity.OfflineActivityResult
	if err = d.db.Model(&record).Select("min(state) as state").Where("activity_id=?", activityID).Find(&record).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	var db = d.db.Table(offlineactivity.TableOfflineActivityInfo).
		Where("id=? and state<10 and state<?", activityID, record.State).
		Update("state", record.State)
	if err = db.Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}
	affectedRow = db.RowsAffected
	if affectedRow > 0 {
		log.Info("update activity state, id=%d, newstate=%d", activityID, record.State)
	} else {
		log.Info("no need to update activity, id=%d, newstate=%d", activityID, record.State)
	}
	return
}
