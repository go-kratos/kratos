package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//ScanSalaryVideoCoupon scan all vip user to salary video coupon.
func (s *Service) ScanSalaryVideoCoupon(c context.Context) (err error) {
	var (
		userInfos  []*model.VipInfoDB
		size       = 100
		endID      int
		now        = time.Now()
		dv         = now.Format("2006_01")
		y          = now.Year()
		m          = now.Month()
		salaryDate = time.Date(y, m, s.c.Property.SalaryDay, 0, 0, 0, 0, time.Local)
	)
	for {
		if endID, err = s.dao.SelUserInfoMaxID(context.TODO()); err != nil {
			log.Error("s.dao.SelMaxID error(%v)", err)
			time.Sleep(time.Minute * 2)
			continue
		}
		break
	}
	page := endID / size
	if endID%size != 0 {
		page++
	}
	for i := 0; i < page; {
		log.Info("salary page(%d) total(%d) ....................................", i, page)
		startID := i * size
		eID := (i + 1) * size
		if userInfos, err = s.dao.SelEffectiveScopeVipList(context.TODO(), startID, eID); err != nil {
			log.Error("s.dao.SelEffectiveScopeVipList error(%v)", err)
			time.Sleep(time.Second * 5)
			continue
		}
		i++
		for _, v := range userInfos {
			time.Sleep(time.Duration(s.c.Property.SalaryVideoCouponnIterval))
			var (
				vipType = model.NotVip
			)
			if v.Status != model.VipStatusNotOverTime && v.Status != model.VipStatusFrozen {
				continue
			}
			if salaryDate.Before(v.OverdueTime.Time()) {
				vipType = model.Vip
				if salaryDate.Before(v.AnnualVipOverdueTime.Time()) {
					vipType = model.AnnualVip
				}
			}
			if vipType == model.NotVip {
				continue
			}
			day := v.OverdueTime.Time().Sub(v.RecentTime.Time()).Hours() / model.DayOfHour
			if day < model.VipDaysMonth {
				continue
			}
			if err = s.salaryCoupon(c, v.Mid, model.TimingSalaryType, int8(vipType), dv, model.CouponSalaryTiming); err != nil {
				err = errors.Wrapf(err, "salaryCoupon mid(%d)(%v)", v.Mid, v)
				log.Error("%+v", err)
				continue
			}
			log.Info("salary suc mid(%d)  ....................................", v.Mid)
		}
	}
	return
}

// salaryCoupon salary coupon.
func (s *Service) salaryCoupon(c context.Context, mid int64, salaryType int8, vipType int8, dv string, atonce int8) (err error) {
	var (
		logs []*model.VideoCouponSalaryLog
		hs   = map[int8]int64{} //  key:coupontype value:salarycount
		ms   map[string]int64   // key:viptype value:salarycount
	)
	if logs, err = s.dao.SalaryVideoCouponList(c, mid, dv); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range logs {
		hs[v.CouponType] = hs[v.CouponType] + v.CouponCount
	}
	for _, v := range s.c.Property.SalaryCouponTypes {
		ms = s.c.Property.SalaryCouponMaps[fmt.Sprintf("%d", v)]
		if len(ms) != 0 {
			if salaryType == model.VipSupplyType {
				if hs[v] == 0 {
					hs[v] = ms[fmt.Sprintf("%d", model.AnnualVip)] - ms[fmt.Sprintf("%d", model.Vip)]
				} else {
					hs[v] = ms[fmt.Sprintf("%d", model.AnnualVip)] - hs[v]
				}
			} else {
				hs[v] = ms[fmt.Sprintf("%d", vipType)] - hs[v]
			}
		}
	}
	for k, count := range hs {
		var (
			token    string
			tokenfmt string
		)
		if count <= 0 {
			continue
		}
		tokenfmt = s.c.Property.SalaryCouponBatchNoMaps[fmt.Sprintf("%d", k)]
		if len(tokenfmt) == 0 {
			continue
		}
		token = fmt.Sprintf(tokenfmt, atonce, dv)
		if err = s.dao.SalaryCoupon(c, mid, k, count, token); err != nil {
			err = errors.Wrapf(err, "s.dao.SalaryCoupon(%d)", mid)
			return
		}
		l := &model.VideoCouponSalaryLog{
			Mid:         mid,
			CouponCount: count,
			State:       model.HadSalaryState,
			Type:        salaryType,
			CouponType:  k,
		}
		if err = s.dao.AddSalaryLog(c, l, dv); err != nil {
			err = errors.WithStack(err)
			return
		}
		if s.c.Property.MsgOpen {
			var (
				title   string
				content string
			)
			title = s.c.Property.SalaryCouponMsgTitleMaps[fmt.Sprintf("%d", k)]
			if len(title) == 0 {
				continue
			}
			if salaryType == model.VipSupplyType {
				content = s.c.Property.SalaryCouponMsgSupplyContentMaps[fmt.Sprintf("%d", k)]
				if len(content) == 0 {
					continue
				}
				content = fmt.Sprintf(content, count)
			} else {
				content = s.c.Property.SalaryCouponMsgContentMaps[fmt.Sprintf("%d", k)]
				if len(content) == 0 {
					continue
				}
			}
			s.sendmessage(func() {
				s.dao.SendMultipMsg(context.TODO(), fmt.Sprintf("%d", mid), content,
					title, model.MsgCouponSalaryMc, model.MsgSystemNotify)
			})
		}
	}
	return
}

// SalaryVideoCouponAtOnce salary video coupon at once.
func (s *Service) SalaryVideoCouponAtOnce(c context.Context, nvip *model.VipUserInfoMsg, ovip *model.VipUserInfoMsg, act string) (res int, err error) {
	if act == _insertAction {
		if err = s.salaryInsertAct(c, nvip); err != nil {
			err = errors.Wrapf(err, "salaryInsertAct (%v)", nvip)
			return
		}
	} else if act == _updateAction {
		if err = s.salaryUpdateAct(c, nvip, ovip); err != nil {
			err = errors.Wrapf(err, "salaryInsertAct (%v)", nvip)
			return
		}
	}
	return
}

func (s *Service) salaryInsertAct(c context.Context, nvip *model.VipUserInfoMsg) (err error) {
	var (
		now        = time.Now()
		otime      time.Time
		aotime     time.Time
		vipType    = model.NotVip
		zeroTime   = now.AddDate(-10, 0, 0)
		salaryType int8
		dv         = now.Format("2006_01")
	)
	otime, err = time.ParseInLocation(model.TimeFormatSec, nvip.OverdueTime, time.Local)
	if err != nil {
		log.Error("time.ParseInLocation error(%v)", errors.Wrapf(err, "time(%s)", nvip.OverdueTime))
		otime = zeroTime
		err = nil
	}
	aotime, err = time.ParseInLocation(model.TimeFormatSec, nvip.AnnualVipOverdueTime, time.Local)
	if err != nil {
		aotime = zeroTime
		err = nil
	}
	if nvip.Status != model.VipStatusNotOverTime && nvip.Status != model.VipStatusFrozen {
		return
	}
	days := otime.Sub(now).Hours() / model.DayOfHour
	if days < model.VipDaysMonth {
		log.Info("cur user not enough send coupon (%+v)", nvip)
		return
	}
	if now.Before(otime) {
		vipType = model.Vip
		if now.Before(aotime) {
			vipType = model.AnnualVip
		}
	}
	switch vipType {
	case model.Vip:
		salaryType = model.NormalVipSalaryType
	case model.AnnualVip:
		salaryType = model.AnnualVipSalaryType
	default:
		return
	}
	if err = s.salaryCoupon(c, int64(nvip.Mid), salaryType, int8(vipType), dv, model.CouponSalaryAtonce); err != nil {
		err = errors.Wrapf(err, "salaryCoupon mid(%d)(%v)", nvip.Mid, nvip)
	}
	return
}

func (s *Service) salaryUpdateAct(c context.Context, nvip *model.VipUserInfoMsg, ovip *model.VipUserInfoMsg) (err error) {
	var (
		ovType     int
		nvType     int
		expire     bool
		now        = time.Now()
		zeroTime   = now.AddDate(-10, 0, 0)
		ntime      time.Time
		oatime     time.Time
		natime     time.Time
		salaryType int8
		dv         = now.Format("2006_01")
	)
	ntime, err = time.ParseInLocation(model.TimeFormatSec, nvip.OverdueTime, time.Local)
	if err != nil {
		log.Error("time.ParseInLocation error(%v)", errors.Wrapf(err, "time(%s)", nvip.OverdueTime))
		ntime = zeroTime
		err = nil
	}
	natime, err = time.ParseInLocation(model.TimeFormatSec, nvip.AnnualVipOverdueTime, time.Local)
	if err != nil {
		natime = zeroTime
		err = nil
	}
	// check OverdueTime time.
	if ntime.Before(now) {
		return
	}
	nvType = model.Vip
	// check AnnualVipOverdueTime time.
	if now.Before(natime) {
		nvType = model.AnnualVip
	}
	// check old vip info expire.
	expire, _ = s.judgeVipExpire(c, ovip)
	if expire {
		//check open days is enough 31
		days := ntime.Sub(now).Hours() / model.DayOfHour
		if days < model.VipDaysMonth {
			log.Info("cur user not enough send coupon (%+v)", nvip)
			return
		}
		if nvType == model.Vip {
			// expire vip -> vip
			salaryType = model.NormalVipSalaryType
		} else if nvType == model.AnnualVip {
			// expire vip -> annual vip
			salaryType = model.AnnualVipSalaryType
		}
	} else {
		if ovip.Type == model.Vip {
			ovType = model.Vip
		}

		oatime, err = time.ParseInLocation(model.TimeFormatSec, ovip.AnnualVipOverdueTime, time.Local)
		if err != nil {
			oatime = zeroTime
			err = nil
		}
		if ovip.Type == model.AnnualVip && oatime.After(now) {
			ovType = model.AnnualVip
		}
		if ovType == model.Vip && nvType == model.AnnualVip {
			// normal vip -> annual vip
			salaryType = model.VipSupplyType
		}
		// short vip -> normal vip
		recentTime := parseTime(ovip.RecentTime)
		otime := parseTime(ovip.OverdueTime)
		days := otime.Sub(recentTime).Hours() / model.DayOfHour
		if days < model.VipDaysMonth {
			//check open days is enough 31
			days := ntime.Sub(now).Hours() / model.DayOfHour
			if days < model.VipDaysMonth {
				log.Info("cur user not enough send coupon (%+v)", nvip)
				return
			}
			if nvType == model.Vip {
				// expire vip -> vip
				salaryType = model.NormalVipSalaryType
			} else if nvType == model.AnnualVip {
				// expire vip -> annual vip
				salaryType = model.AnnualVipSalaryType
			}
		}

	}
	switch salaryType {
	case model.NormalVipSalaryType, model.AnnualVipSalaryType, model.VipSupplyType:
		if err = s.salaryCoupon(c, int64(nvip.Mid), salaryType, int8(nvType), dv, model.CouponSalaryAtonce); err != nil {
			err = errors.Wrapf(err, "salaryCoupon mid(%d)(%v)", int64(nvip.Mid), nvip)
		}
	}
	return
}

// judgeVipExpire judge vip is expire.
func (s *Service) judgeVipExpire(c context.Context, v *model.VipUserInfoMsg) (expire bool, err error) {
	var (
		now         = time.Now()
		overdueTime time.Time
		zeroTime    = now.AddDate(-10, 0, 0)
	)
	if v.Status != model.VipStatusNotOverTime && v.Status != model.VipStatusFrozen {
		expire = true
		return
	}
	overdueTime, err = time.ParseInLocation(model.TimeFormatSec, v.OverdueTime, time.Local)
	if err != nil {
		log.Error("time.ParseInLocation error(%v)", errors.Wrapf(err, "time(%s)", v.OverdueTime))
		overdueTime = zeroTime
		err = nil
	}
	if overdueTime.Before(now) {
		expire = true
		return
	}
	return
}
