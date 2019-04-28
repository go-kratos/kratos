package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	xtime "go-common/library/time"
)

func changeType(hv *model.HandlerVip) (r int) {
	if hv.OldVipUser == nil {
		r = model.VipChangeOpen
	} else {
		r = model.VipChangeModify
	}
	return
}

func getOpenVipMsg(cs int, d int64) (r string) {
	var (
		s string
	)
	if d/model.VipDaysMonth != 0 && d%model.VipDaysMonth == 0 {
		s = strconv.FormatInt(int64(d/model.VipDaysMonth), 10) + "个月"
	} else {
		s = strconv.FormatInt(int64(d), 10) + "天"
	}
	if model.VipChangeOpen == cs {

		r = fmt.Sprintf(model.VipOpenMsg, s)
	} else {
		r = fmt.Sprintf(model.VipOpenKMsg, s)
	}
	return
}

//UpdateVipWithHistory .
func (s *Service) UpdateVipWithHistory(c context.Context, tx *sql.Tx, v *model.VipChangeBo) (r *model.HandlerVip, err error) {
	var (
		days      int64
		vch       = new(model.VipChangeHistory)
		historyID int64
	)
	if v.Months == 0 {
		days = v.Days
	} else {
		year := v.Months / 12
		month := v.Months % 12
		days = int64(int(year)*model.VipDaysYear + int(month)*model.VipDaysMonth)
	}

	if r, err = s.vipDuration(c, tx, v.Mid, days); err != nil {
		return
	}
	vch.BatchCodeID = v.BatchCodeID
	vch.BatchID = v.BatchID
	vch.ChangeTime = xtime.Time(time.Now().Unix())
	vch.ChangeType = v.ChangeType
	vch.Days = days
	vch.Mid = v.Mid
	vch.RelationID = v.RelationID
	vch.Remark = v.Remark
	vch.OperatorID = v.OperatorID
	if historyID, err = s.saveVipHistory(tx, vch); err != nil {
		return
	}
	r.HistoryID = historyID
	r.Days = days
	return
}

func (s *Service) saveVipHistory(tx *sql.Tx, v *model.VipChangeHistory) (id int64, err error) {
	id, err = s.dao.InsertVipChangeHistory(tx, v)
	return
}

func (s *Service) vipDuration(c context.Context, tx *sql.Tx, mid int64, days int64) (r *model.HandlerVip, err error) {
	var (
		ou *model.VipInfoDB
		re *model.VipInfoDB
	)
	r = new(model.HandlerVip)
	if days < 1 {
		err = ecode.VipDaysErr
	}
	if ou, err = s.dao.TxSelVipUserInfo(tx, mid); err != nil {
		return
	}
	if ou == nil {
		if re, err = s.insertVipUser(c, tx, mid, days); err != nil {
			return
		}
	} else {
		if re, err = s.updateVipUser(c, tx, mid, days, ou); err != nil {
			return
		}
	}
	r.OldVipUser = ou
	r.VipUser = re
	r.Mid = mid
	return
}

func (s *Service) insertVipUser(c context.Context, tx *sql.Tx, mid int64, days int64) (r *model.VipInfoDB, err error) {
	var (
		curTime time.Time
		it      time.Time
	)
	r = new(model.VipInfoDB)
	if days < 1 {
		err = ecode.VipDaysErr
	}
	if it, err = time.ParseInLocation("2006-01-02", "2016-01-01", time.Local); err != nil {
		return
	}
	r.AnnualVipOverdueTime = xtime.Time(it.Unix())
	if curTime, err = time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local); err != nil {
		return
	}
	curTime = curTime.AddDate(0, 0, 1)
	if days > model.VipDaysYear {
		r.VipType = model.AnnualVip
		curTime = curTime.AddDate(0, 0, int(days))
		r.VipOverdueTime = xtime.Time(curTime.Unix())
		r.AnnualVipOverdueTime = r.VipOverdueTime
	} else {
		curTime = curTime.AddDate(0, 0, int(days))
		r.VipType = model.Vip
		r.VipOverdueTime = xtime.Time(curTime.Unix())
	}
	ct := xtime.Time(time.Now().Unix())
	r.Mid = int64(mid)
	r.VipStatus = model.NotExpired
	r.VipStartTime = ct
	r.VipRecentTime = ct
	err = s.dao.InsertVipUserInfo(tx, r)
	return
}

func (s *Service) updateVipUser(c context.Context, tx *sql.Tx, mid int64, days int64, ou *model.VipInfoDB) (r *model.VipInfoDB, err error) {
	var a int64
	r = handlerVipTime(days, ou)
	r.VipStatus = model.NotExpired
	r.Ver = ou.Ver + 1
	r.ID = ou.ID
	if a, err = s.dao.UpdateVipUserInfo(tx, r, ou.Ver); err != nil {
		return
	}
	if a <= 0 {
		err = ecode.VipUpdateErr
	}
	return
}

func handlerVipTime(days int64, ou *model.VipInfoDB) (r *model.VipInfoDB) {
	var (
		ct time.Time
		rt time.Time
	)
	r = new(model.VipInfoDB)
	ct, _ = time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	r.AnnualVipOverdueTime = ou.AnnualVipOverdueTime
	if ct.Before(ou.VipOverdueTime.Time()) {
		rt = ou.VipOverdueTime.Time()
	} else {
		rt = ct.Add(time.Hour * 24)
	}
	rt = rt.AddDate(0, 0, int(days))

	r.VipOverdueTime = xtime.Time(rt.Unix())
	if (rt.Sub(ct).Hours() / 24) >= model.VipDaysYear {
		r.VipType = model.AnnualVip
		r.AnnualVipOverdueTime = r.VipOverdueTime
	} else {
		if ou.VipType == model.AnnualVip {
			r.VipType = model.AnnualVip
		} else {
			r.VipType = model.Vip
		}
	}
	curTime := time.Now()
	if ou.VipStatus == model.Expire || ou.VipOverdueTime.Time().Before(curTime) {
		r.VipRecentTime = xtime.Time(curTime.Unix())
	} else {
		r.VipRecentTime = ou.VipRecentTime
	}

	return
}
