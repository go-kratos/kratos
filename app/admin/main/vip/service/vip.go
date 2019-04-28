package service

import (
	"context"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// Drawback Drawback
func (s *Service) Drawback(c context.Context, days int, mid int64, usename, remark string) (err error) {
	var (
		userInfo *model.VipUserInfo
		history  *model.VipChangeHistory
		ot       time.Time
		tx       *sql.Tx
		affected int64
	)
	if userInfo, err = s.dao.SelVipUserInfo(c, mid); err != nil {
		return
	}
	if userInfo == nil {
		err = ecode.VipUserInfoNotExit
		return
	}
	ct := time.Now()
	overdueTIme := userInfo.VipOverdueTime
	ot = overdueTIme.Time().AddDate(0, 0, -days)
	if ot.Before(ct) {
		ot = ct
		userInfo.VipStatus = model.VipStatusOverTime
	}
	userInfo.VipOverdueTime = xtime.Time(ot.Unix())
	if userInfo.AnnualVipOverdueTime.Time().After(ot) && userInfo.VipType == model.AnnualVip {
		userInfo.VipType = model.Vip
		userInfo.AnnualVipOverdueTime = xtime.Time(ot.Unix())
	}
	history = new(model.VipChangeHistory)
	history.Mid = mid
	history.ChangeTime = xtime.Time(time.Now().Unix())
	history.Days = -days
	history.ChangeType = model.ChangeTypeSystemDrawback
	history.OperatorID = usename
	history.Remark = remark

	if tx, err = s.dao.BeginTran(context.TODO()); err != nil {
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
	}()
	if affected, err = s.dao.UpdateVipUserInfo(tx, userInfo); err != nil {
		return
	}
	if affected > 0 {
		if _, err = s.dao.InsertVipChangeHistory(tx, history); err != nil {
			return
		}
		if err = s.dao.DelBcoinSalary(tx, mid, userInfo.VipOverdueTime); err != nil {
			return
		}
	}
	return
}

// ExchangeVipByDays exchange vip days
//func (s *Service) exchangeVip(c context.Context, tx *sql.Tx, mid, batchID, unit int, remark, username string) (hv *inModel.HandlerVip, err error) {
//	var r = new(inModel.VipChangeBo)
//	if mid == 0 {
//		err = ecode.VipMidErr
//		return
//	}
//	if len(remark) == 0 || len(remark) > 200 {
//		err = ecode.VipRemarkErr
//		return
//	}
//	r.BatchID = int64(batchID)
//	r.ChangeType = model.ChangeTypeSystem
//	r.Remark = remark
//	r.OperatorID = username
//	r.Days = int32(unit)
//	r.Mid = int64(mid)
//	log.Info("%+v ", r)
//	if hv, err = s.vipRPC.UpdateVipWithHistory(context.TODO(), r); err != nil {
//		fmt.Printf("rpc %+v err(%+v) \n", hv, err)
//		log.Error("rpc %+v err(%+v) ", hv, err)
//		return
//	}
//	return
//}

// HistoryPage history page.
func (s *Service) HistoryPage(c context.Context, u *model.UserChangeHistoryReq) (res []*model.VipChangeHistory, count int, err error) {
	if count, err = s.dao.HistoryCount(c, u); err != nil {
		log.Error("%+v", err)
		return
	}
	if count == 0 {
		return
	}
	if res, err = s.dao.HistoryList(c, u); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

// VipInfo .
func (s *Service) VipInfo(c context.Context, mid int64) (res *model.VipUserInfo, err error) {
	if res, err = s.dao.SelVipUserInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
