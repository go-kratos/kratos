package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"
)

const (
	memo = "年费会员发B券"
)

// OldProcesserHandler old processer handler.
func (s *Service) OldProcesserHandler(c context.Context, hv *model.OldHandlerVip, ip string) (err error) {
	if hv == nil || hv.VipUser == nil {
		return
	}
	go s.dao.SendMultipMsg(c, fmt.Sprintf("%d", hv.Mid), oldGetOpenVipMsg(oldChangeType(hv), hv.Days), model.VipOpenMsgTitle, model.VipOpenMsgCode, ip, model.VipSystemNotify)
	err = s.oldAddBcoin(c, hv.Mid, hv.VipUser, ip)
	return
}

func (s *Service) oldAddBcoin(c context.Context, mid int64, vu *model.VipUserInfo, ip string) (err error) {
	var (
		ct  time.Time
		b   time.Time
		vbs *model.VipBcoinSalary
	)
	if vu != nil && model.AnnualVip != vu.VipType {
		return
	}

	b = time.Now()
	ct = time.Now()
	ct = ct.Add(time.Duration((1 - ct.Day())) * 24 * time.Hour)
	ct, _ = time.Parse("2006-01-02", ct.Format("2006-01-02"))
	if vbs, err = s.dao.OldSelLastBcoin(c, mid); err != nil {
		return
	}
	if vbs != nil && vbs.Month.Time().Unix() > ct.Unix() {
		b = vbs.Month.Time()
	} else if vu.AnnualVipOverdueTime.Time().After(time.Now()) {
		if err = s.oldSendBcoinNow(c, mid, ip); err != nil {
			return
		}
	}
	err = s.oldSetVipBcoin(c, mid, b, vu.AnnualVipOverdueTime.Time())
	return
}

func (s *Service) oldSendBcoinNow(c context.Context, mid int64, ip string) (err error) {
	var (
		bb  *model.BcoinSendBo
		vbs = new(model.VipBcoinSalary)
		vc  = s.vipConfig[model.AnnualVipBcoinCouponActivityID]
	)
	if bb, err = s.oldSendInfo(); err != nil {
		return
	}
	if vc == nil {
		err = ecode.VipConfigNotExitErr
		return
	}

	//接入支付中心
	s.dao.SendBcoinCoupon(c, strconv.FormatInt(int64(mid), 10), vc.Content, bb.Amount, bb.DueDate.Time())

	vbs.Amount = bb.Amount
	vbs.Mid = mid
	vbs.GiveNowStatus = 1
	vbs.Memo = memo
	vbs.Month = xtime.Time(time.Now().Unix())
	if err = s.dao.OldInsertVipBcoinSalary(c, vbs); err != nil {
		return
	}
	s.dao.SendMultipMsg(c, strconv.FormatInt(int64(mid), 10), fmt.Sprintf(model.VipBcoinGiveContext, bb.Amount, bb.DayOfMonth), model.VipBcoinGiveTitle, model.VipBcoinGiveMsgCode, ip, model.VipSystemNotify)
	return
}

func (s *Service) oldSetVipBcoin(c context.Context, mid int64, b, e time.Time) (err error) {
	var (
		bb  *model.BcoinSendBo
		vbs = new(model.VipBcoinSalary)
	)
	if bb, err = s.oldSendInfo(); err != nil {
		return
	}
	vbs.Amount = bb.Amount
	vbs.Status = 0
	vbs.Mid = mid
	vbs.Memo = memo
	vbs.GiveNowStatus = 0

	b = b.AddDate(0, 1, int(bb.DayOfMonth)-b.Day())
	e = e.AddDate(0, 1, 1-e.Day())
	for b.Unix() < e.Unix() {
		vbs.Month = xtime.Time(b.Unix())
		if err = s.dao.OldInsertVipBcoinSalary(c, vbs); err != nil {
			return
		}
		b = b.AddDate(0, 1, 0)
	}
	return
}

func (s *Service) oldSendInfo() (r *model.BcoinSendBo, err error) {
	var (
		c      time.Time
		d      = s.vipConfig[model.AnnualVipBcoinDay]
		a      = s.vipConfig[model.AnnualVipBcoinCouponMoney]
		day    int64
		amount int64
	)
	r = new(model.BcoinSendBo)
	if d == nil {
		err = ecode.VipConfigNotExitErr
		return
	}
	if a == nil {
		err = ecode.VipConfigNotExitErr
		return
	}
	if day, err = strconv.ParseInt(d.Content, 10, 64); err != nil {
		return
	}
	if amount, err = strconv.ParseInt(a.Content, 10, 64); err != nil {
		return
	}
	r.Amount = amount
	r.DayOfMonth = day
	c = time.Now()
	c = c.AddDate(0, 1, int(day)-1-c.Day())
	r.DueDate = xtime.Time(c.Unix())
	return
}
