package service

import (
	"github.com/jinzhu/gorm"
	"go-common/app/job/main/up/dao/upcrm"
	"go-common/app/job/main/up/model/signmodel"
	"go-common/app/service/main/upcredit/mathutil"
	"go-common/library/log"
	"time"
)

// change a lot state for sign up
// state, 签约状态，正常或过期
// due_wan, 签约状态即将过期
// pay_expire_state，付费状态

//CheckStateJob 检查sign_up中的状态
func (s *Service) CheckStateJob(date time.Time) {
	log.Info("start run state job, date=%s", date)
	s.checkState(date)
	log.Info("finish run state job, date=%s", date)
}

func (s *Service) checkState(date time.Time) {
	// 更新签约状态：
	var crmdb = s.crmdb.GetDb()
	// 签约30天内到期
	var signExpireList []*signmodel.SignUp
	var signDueList []*signmodel.SignUp
	var limit, offset = 200, 0
	var count = limit
	var err error
	var signEndDate = date.AddDate(0, 0, 30)
	for count == limit {
		var items []*signmodel.SignUp
		// 把endDate在[-60, +30]范围内的用户全都找出来，过一遍
		err = crmdb.Offset(offset).Limit(limit).Where("end_date<? and end_date>?", signEndDate, date.AddDate(0, -4, 0)).Find(&items).Error
		if err != nil {
			log.Error("fail to get sign ups, err=%v", err)
			break
		}
		count = len(items)
		offset += count
		for _, v := range items {
			if v.EndDate.Time().Before(date) && v.State == signmodel.SignStateOnSign {
				// 过期
				signExpireList = append(signExpireList, v)
			} else if !(v.EndDate.Time().Before(date) || v.EndDate.Time().After(signEndDate)) && v.DueWarn <= signmodel.DueWarnNoWarn {
				// 快到期
				signDueList = append(signDueList, v)
			}
		}
	}

	var expireIDList []uint32
	for _, v := range signExpireList {
		expireIDList = append(expireIDList, v.ID)
	}

	err = updateListWithLimit(crmdb, "state", signmodel.SignStateExpire, expireIDList, limit)
	if err != nil {
		log.Error("fail to update sign state, err=%v", err)
		return
	}
	var dueIDList []uint32
	for _, v := range signDueList {
		dueIDList = append(dueIDList, v.ID)
	}

	err = updateListWithLimit(crmdb, "due_warn", signmodel.DueWarnWarn, dueIDList, limit)
	if err != nil {
		log.Error("fail to update sign due_warn, err=%v", err)
		return
	}
	// 付款7天内到期
	offset = 0
	count = limit
	var payDueDate = date.AddDate(0, 0, 7)
	var payExpireMap = map[uint32]int8{}
	for count == limit {
		var items []*signmodel.SignPay
		// 把endDate在[-max, +7]范围内的用户全都找出来，过一遍
		err = crmdb.Offset(offset).Limit(limit).Where("due_date<?", payDueDate).Find(&items).Error
		if err != nil {
			log.Error("fail to get sign ups, err=%v", err)
			break
		}
		count = len(items)
		offset += count
		for _, v := range items {
			// 如果有即将到期，则立即标记为即将到期
			if v.State == upcrm.PayStateUnpay {
				payExpireMap[v.SignID] = signmodel.PayExpireStateDue
			} else {
				// 如果没有即将到期的状态，则把其标记为未到期
				// 即：只有所有的付款都已标记完成，才会认为是未到期状态
				if payExpireMap[v.SignID] != signmodel.PayExpireStateDue {
					payExpireMap[v.SignID] = signmodel.PayExpireStateNormal
				}
			}
		}
	}
	// 更新到期状态
	var stateSignIDListMap = map[int8][]uint32{}
	for k, v := range payExpireMap {
		stateSignIDListMap[v] = append(stateSignIDListMap[v], k)
	}
	// 分区更新
	limit = 200
	for state, list := range stateSignIDListMap {
		err = updateListWithLimit(crmdb, "pay_expire_state", state, list, limit)
		if err != nil {
			log.Error("err update pay_expire_state, err=%v", err)
			return
		}
	}
}

func updateListWithLimit(crmdb *gorm.DB, field string, state int8, list []uint32, limit int) (err error) {
	for begin := 0; begin < len(list); begin += limit {
		var end = mathutil.Min(begin+limit, len(list))
		var needUpdate = list[begin:end]
		err = updateSignTable(crmdb, needUpdate, field, state)
		if err != nil {
			log.Error("fail to update state [%s], err=%v", field, err)
			return
		}
	}
	return
}
func updateSignTable(crmdb *gorm.DB, ids interface{}, field string, state int8) (err error) {
	err = crmdb.Table(signmodel.TableNameSignUp).Where("id in (?)", ids).Update(field, state).Error
	return
}
