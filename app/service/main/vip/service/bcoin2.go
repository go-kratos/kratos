package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

//ProcesserHandler def.
func (s *Service) ProcesserHandler(c context.Context, hv *model.HandlerVip, ip string) (err error) {
	if hv.ToMid == 0 {
		go s.dao.SendMultipMsg(c, strconv.FormatInt(int64(hv.Mid), 10), getOpenVipMsg(changeType(hv), hv.Days), model.VipOpenMsgTitle, model.VipOpenMsgCode, ip, model.VipSystemNotify)
	} else {
		go func() {
			var (
				zm   *memmdl.BaseInfo
				bm   *memmdl.BaseInfo
				err1 error
			)
			if zm, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: hv.Mid}); err1 != nil {
				log.Error("s.accRPC.Info3(%+v) error(%+v)", hv, err1)
				return
			}
			if bm, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: hv.ToMid}); err1 != nil {
				log.Error("s.accRPC.Info3(%+v) error(%+v)", hv, err1)
				return
			}
			if zm == nil || bm == nil {
				return
			}
			year := hv.Days / model.VipDaysYear
			sulperDays := hv.Days % model.VipDaysYear
			months := year/12 + sulperDays/model.VipDaysMonth
			//赠送人
			s.dao.SendMultipMsg(c, fmt.Sprintf("%d", hv.Mid),
				fmt.Sprintf(_vipZengsongMsgContext, bm.Name, hv.ToMid, months),
				_vipZengsongMsgTitle,
				model.VipCustomizeMsgCode,
				ip,
				model.VipSystemNotify)
			//收礼物人
			s.dao.SendMultipMsg(c, fmt.Sprintf("%d", hv.ToMid),
				fmt.Sprintf(_vipShouliwuMsgContext, zm.Name, hv.Mid, months, hv.Mid, months, time.Now().Unix(), hv.Mid, months, time.Now().Unix()),
				_vipShouliwuMsgtitle,
				model.VipCustomizeMsgCode,
				ip,
				model.VipSystemNotify)
		}()

	}

	err = s.addBcoin(c, hv.Mid, hv.VipUser, ip)
	return
}

func (s *Service) addBcoin(c context.Context, mid int64, vu *model.VipInfoDB, ip string) (err error) {
	var (
		ct  time.Time
		b   time.Time
		vbs *model.VipBcoinSalary
	)
	if model.AnnualVip != vu.VipType {
		return
	}

	b = time.Now()
	ct = time.Now()
	ct = ct.Add(time.Duration((1 - ct.Day())) * 24 * time.Hour)
	ct, _ = time.Parse("2006-01-02", ct.Format("2006-01-02"))
	if vbs, err = s.dao.SelLastBcoin(c, mid); err != nil {
		return
	}
	if vbs != nil && vbs.Month.Time().Unix() > ct.Unix() {
		b = vbs.Month.Time()
	} else {
		if err = s.sendBcoinNow(c, mid, ip); err != nil {
			return
		}
	}
	err = s.setVipBcoin(c, mid, b, vu.AnnualVipOverdueTime.Time())
	return
}

func (s *Service) sendBcoinNow(c context.Context, mid int64, ip string) (err error) {
	var (
		bb  *model.BcoinSendBo
		vbs = new(model.VipBcoinSalary)
		vc  = s.vipConfig[model.AnnualVipBcoinCouponActivityID]
	)
	if bb, err = s.sendInfo(); err != nil {
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
	vbs.Memo = "年费会员发B券"
	vbs.Month = xtime.Time(time.Now().Unix())
	if err = s.dao.InsertVipBcoinSalary(c, vbs); err != nil {
		return
	}
	s.dao.SendMultipMsg(c, strconv.FormatInt(int64(mid), 10), fmt.Sprintf(model.VipBcoinGiveContext, bb.Amount, bb.DayOfMonth), model.VipBcoinGiveTitle, model.VipBcoinGiveMsgCode, ip, model.VipSystemNotify)
	return
}

func (s *Service) setVipBcoin(c context.Context, mid int64, b, e time.Time) (err error) {
	var (
		bb  *model.BcoinSendBo
		vbs = new(model.VipBcoinSalary)
	)
	if bb, err = s.sendInfo(); err != nil {
		return
	}
	vbs.Amount = bb.Amount
	vbs.Status = 0
	vbs.Mid = mid
	vbs.Memo = "年费会员发放B币"
	vbs.GiveNowStatus = 0

	b = b.AddDate(0, 1, int(bb.DayOfMonth)-b.Day())
	e = e.AddDate(0, 1, 1-e.Day())
	for b.Unix() < e.Unix() {
		vbs.Month = xtime.Time(b.Unix())
		if err = s.dao.InsertVipBcoinSalary(c, vbs); err != nil {
			return
		}
		b = b.AddDate(0, 1, 0)
	}
	return
}

func (s *Service) sendInfo() (r *model.BcoinSendBo, err error) {
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
	if day, err = strconv.ParseInt(d.Content, 10, 0); err != nil {
		return
	}
	if amount, err = strconv.ParseInt(a.Content, 10, 0); err != nil {
		return
	}
	r.Amount = amount
	r.DayOfMonth = day
	c = time.Now()
	c = c.AddDate(0, 1, int(day)-1-c.Day())
	r.DueDate = xtime.Time(c.Unix())
	return
}
