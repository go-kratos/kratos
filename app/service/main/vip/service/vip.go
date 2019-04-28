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

func (s *Service) oldCleanCacheAndNotify(c context.Context, hv *model.OldHandlerVip, ip string) (err error) {
	if hv == nil {
		return
	}
	t := 0
	if hv.OldVipUser != nil {
		t = 1
	}
	err = s.dao.SendCleanCache(c, hv.Mid, hv.Months, hv.Days, t, ip)
	return
}

func oldChangeType(hv *model.OldHandlerVip) (r int) {
	if hv.OldVipUser == nil {
		r = model.VipChangeOpen
	} else {
		r = model.VipChangeModify
	}
	return
}

func oldGetOpenVipMsg(cs int, d int64) (r string) {
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

//OldUpdateVipWithHistory update Vip history.
func (s *Service) OldUpdateVipWithHistory(c context.Context, tx *sql.Tx, v *model.VipChangeBo) (r *model.OldHandlerVip, err error) {
	var (
		vch       = new(model.VipChangeHistory)
		historyID int64
	)
	if r, err = s.oldVipDuration(c, tx, v.Mid, v.Days); err != nil {
		return
	}

	vch.BatchID = v.BatchID
	vch.ChangeTime = xtime.Time(time.Now().Unix())
	vch.ChangeType = v.ChangeType
	vch.Days = v.Days
	vch.Mid = v.Mid
	vch.RelationID = v.RelationID
	vch.Remark = v.Remark
	if historyID, err = s.oldSaveVipHistory(c, tx, vch); err != nil {
		return
	}
	r.HistoryID = historyID
	r.Days = v.Days
	return
}

func (s *Service) oldSaveVipHistory(c context.Context, tx *sql.Tx, v *model.VipChangeHistory) (id int64, err error) {
	id, err = s.dao.OldInsertVipChangeHistory(c, tx, v)
	return
}

func (s *Service) oldVipDuration(c context.Context, tx *sql.Tx, mid, days int64) (r *model.OldHandlerVip, err error) {
	var (
		ou *model.VipUserInfo
		re *model.VipUserInfo
	)
	r = new(model.OldHandlerVip)
	if days < 1 {
		err = ecode.VipDaysErr
		return
	}
	if ou, err = s.dao.OldSelVipUserInfo(c, mid); err != nil {
		return
	}
	if ou == nil {
		if re, err = s.oldInsertVipUser(c, tx, mid, days); err != nil {
			return
		}
	} else {
		if re, err = s.oldUpdateVipUser(c, tx, mid, days, ou); err != nil {
			return
		}
	}
	r.OldVipUser = ou
	r.VipUser = re
	r.Mid = mid
	return
}

func (s *Service) oldInsertVipUser(c context.Context, tx *sql.Tx, mid, days int64) (r *model.VipUserInfo, err error) {
	var (
		curTime time.Time
		it      time.Time
	)
	r = new(model.VipUserInfo)
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
	if days >= model.VipDaysYear {
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
	r.Mid = mid
	r.VipStatus = model.VipStatusNotOverTime
	r.VipStartTime = ct
	r.VipRecentTime = ct
	r.Wander = 0
	r.AccessStatus = 0
	err = s.dao.OldInsertVipUserInfo(c, tx, r)
	return
}

func (s *Service) oldUpdateVipUser(c context.Context, tx *sql.Tx, mid, days int64, ou *model.VipUserInfo) (r *model.VipUserInfo, err error) {
	var a int64
	r = oldHandlerVipTime(int(days), ou)
	r.VipStatus = model.VipStatusNotOverTime
	r.ID = ou.ID
	if a, err = s.dao.OldUpdateVipUserInfo(c, tx, r); err != nil {
		return
	}
	if a <= 0 {
		err = ecode.VipUpdateErr
	}
	return
}

func oldHandlerVipTime(days int, ou *model.VipUserInfo) (r *model.VipUserInfo) {
	var (
		ct time.Time
		rt time.Time
	)
	r = new(model.VipUserInfo)
	ct, _ = time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	r.AnnualVipOverdueTime = ou.AnnualVipOverdueTime
	if ct.Before(ou.VipOverdueTime.Time()) {
		rt = ou.VipOverdueTime.Time()
	} else {
		rt = ct.Add(time.Hour * 24)
	}
	rt = rt.AddDate(0, 0, days)
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
	if ou.VipStatus == model.VipStatusOverTime || ou.VipOverdueTime.Time().Before(curTime) {
		r.VipRecentTime = xtime.Time(curTime.Unix())
	} else {
		r.VipRecentTime = ou.VipRecentTime
	}
	return
}
