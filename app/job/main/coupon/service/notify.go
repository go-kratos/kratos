package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// Notify notify.
func (s *Service) Notify(c context.Context, msg *model.MsgCanal) (err error) {
	var (
		mid                     int64
		token                   string
		ok, ok1                 bool
		ct                      int64
		couponToken, batchToken string
	)
	if strings.Contains(msg.Table, _couponTable) {
		if msg.Action != _updateAct {
			return
		}
		if mid, token, ct, ok, err = s.conventMsg(c, msg); err != nil {
			err = errors.WithStack(err)
			return
		}
	} else if msg.Table == _orderTable {
		if msg.Action != _insertAct {
			return
		}
		if mid, token, ct, ok, err = s.conventOrderMsg(c, msg); err != nil {
			err = errors.WithStack(err)
			return
		}
	} else if strings.Contains(msg.Table, _couponAllowanceTable) { // 元旦活动
		if msg.Action != _updateAct {
			return
		}
		if mid, couponToken, batchToken, ok1, err = s.conventAllowanceInfoMsg(c, msg); err != nil {
			err = errors.WithStack(err)
			return
		}
		if ok1 {
			if _, err = s.dao.UpdateUserCard(c, mid, model.Used, couponToken, batchToken); err != nil {
				err = errors.WithStack(err)
				return
			}
			if err = s.dao.DelPrizeCardsKey(c, mid, s.c.NewYearConf.ActID); err != nil {
				err = errors.WithStack(err)
				return
			}
		}
	}
	if !ok {
		return
	}
	arg := &model.NotifyParam{
		Mid:         mid,
		CouponToken: token,
		NotifyURL:   s.c.Properties.BangumiNotifyURL,
		Type:        ct,
	}
	if err = s.CheckCouponDeliver(c, arg); err != nil {
		log.Error("CheckCouponDeliver fail arg(%v) err(%v)", arg, err)
		arg.NotifyCount++
		s.notifyChan <- arg
		return
	}
	return
}

func (s *Service) conventMsg(c context.Context, msg *model.MsgCanal) (mid int64, token string, ct int64, ok bool, err error) {
	ok = true
	cnew := new(model.CouponInfo)
	if err = json.Unmarshal(msg.New, cnew); err != nil {
		err = errors.WithStack(err)
		return
	}
	cold := new(model.CouponInfo)
	if err = json.Unmarshal(msg.Old, cold); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cold.State != model.NotUsed || cnew.State != model.InUse {
		ok = false
	}
	mid = cnew.Mid
	token = cnew.CouponToken
	ct = cnew.CouponType
	return
}

func (s *Service) conventOrderMsg(c context.Context, msg *model.MsgCanal) (mid int64, token string, ct int64, ok bool, err error) {
	ok = true
	cnew := new(model.CouponOrder)
	if err = json.Unmarshal(msg.New, cnew); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cnew.State != model.InPay {
		ok = false
	}
	mid = cnew.Mid
	token = cnew.OrderNo
	ct = int64(cnew.CouponType)
	return
}

func (s *Service) conventAllowanceInfoMsg(c context.Context, msg *model.MsgCanal) (mid int64, couponToken, batchToken string, ok bool, err error) {
	ok = true
	cnew := new(model.CouponAllowanceInfo)
	if err = json.Unmarshal(msg.New, cnew); err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Info("conventAllowanceInfoMsg(%+v)", cnew)
	if cnew.State != model.Used || cnew.AppID != 1 || cnew.Origin != model.AllowanceBusinessNewYear {
		ok = false
	}
	mid = cnew.MID
	couponToken = cnew.CouponToken
	batchToken = cnew.BatchToken
	return
}

//CheckCouponDeliver check coupon deliver
func (s *Service) CheckCouponDeliver(c context.Context, arg *model.NotifyParam) (err error) {
	switch arg.Type {
	case model.BangumiVideo:
		err = s.CouponDeliver(c, arg)
	case model.Cartoon:
		err = s.CouponCartoonDeliver(c, arg)
	}
	return
}

// CouponDeliver def.
func (s *Service) CouponDeliver(c context.Context, arg *model.NotifyParam) (err error) {
	var (
		data   *model.CallBackRet
		cp     *model.CouponInfo
		nstate int8
	)
	if cp, err = s.dao.CouponInfo(c, arg.Mid, arg.CouponToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp == nil {
		log.Warn("notify coupon is nil(%v)", arg)
		return
	}
	if cp.State != model.InUse {
		log.Warn("notify coupon had deal with(%v)", arg)
		return
	}
	if data, err = s.dao.NotifyRet(c, arg.NotifyURL, cp.CouponToken, cp.OrderNO, "127.0.0.1"); err != nil {
		err = errors.WithStack(err)
		return
	}
	if data.Ver == cp.UseVer {
		err = fmt.Errorf("coupon ver not change resp(%v) db(%v)", data, cp)
		return
	}
	switch data.IsPaid {
	case model.PaidSuccess:
		nstate = model.Used
	case model.Unpaid:
		nstate = model.NotUsed
	default:
		log.Warn("state not found resp(%v) db(%v)", data, cp)
		return
	}
	log.Info("update coupon state(%s,%d,%d,%d,%d)", cp.CouponToken, cp.Mid, nstate, data.Ver, cp.Ver)
	if err = s.updateCouponState(c, cp, nstate, data); err != nil {
		log.Error("updateCouponState fail %+v", err)
		return
	}
	return
}

func (s *Service) updateCouponState(c context.Context, cp *model.CouponInfo, nstate int8, data *model.CallBackRet) (err error) {
	var (
		tx  *sql.Tx
		aff int64
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback %+v", err)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit %+v", err)
		}
	}()
	if aff, err = s.dao.UpdateCoupon(c, tx, cp.Mid, nstate, data.Ver, cp.Ver, cp.CouponToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("coupon deal fail (%v) db(%v)", data, cp)
		return
	}
	l := &model.CouponChangeLog{}
	l.CouponToken = cp.CouponToken
	l.Mid = cp.Mid
	l.State = nstate
	l.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.InsertPointHistory(c, tx, l); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.dao.DelCouponsCache(c, cp.Mid, int8(cp.CouponType))
	return
}

// CouponCartoonDeliver coupon cartoon deliver def.
func (s *Service) CouponCartoonDeliver(c context.Context, arg *model.NotifyParam) (err error) {
	var (
		data   *model.CallBackRet
		o      *model.CouponOrder
		nstate int8
	)
	if o, err = s.dao.ByOrderNo(c, arg.CouponToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if o == nil {
		log.Warn("notify coupon order is nil(%v)", arg)
		return
	}
	if o.State != model.InPay {
		log.Warn("notify coupon order had deal with(%v)", arg)
		return
	}
	if data, err = s.dao.NotifyRet(c, arg.NotifyURL, o.OrderNo, o.ThirdTradeNo, "127.0.0.1"); err != nil {
		err = errors.WithStack(err)
		return
	}
	if data.Ver == o.UseVer {
		err = fmt.Errorf("coupon order ver not change resp(%v) db(%v)", data, o)
		return
	}
	switch data.IsPaid {
	case model.PaidSuccess:
		nstate = model.PaySuccess
	case model.Unpaid:
		nstate = model.PayFaild
	default:
		log.Warn("order state not found resp(%v) db(%v)", data, o)
		return
	}
	log.Info("update coupon order state(%s,%d,%d,%d,%d)", o.OrderNo, o.Mid, nstate, data.Ver, o.UseVer)
	if err = s.UpdateOrderState(c, o, nstate, data); err != nil {
		log.Error("updateCouponState fail %+v", err)
		return
	}
	return
}

// UpdateOrderState update order state.
func (s *Service) UpdateOrderState(c context.Context, o *model.CouponOrder, nstate int8, data *model.CallBackRet) (err error) {
	var (
		tx  *sql.Tx
		aff int64
		ls  []*model.CouponBalanceChangeLog
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback %+v", err)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit %+v", err)
		}
	}()
	if aff, err = s.dao.UpdateOrderState(c, tx, o.Mid, nstate, data.Ver, o.Ver, o.OrderNo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("coupon order deal fail (%v) db(%v)", data, o)
		return
	}
	// add order log.
	ol := new(model.CouponOrderLog)
	ol.OrderNo = o.OrderNo
	ol.Mid = o.Mid
	ol.State = nstate
	ol.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddOrderLog(c, tx, ol); err != nil {
		err = errors.WithStack(err)
		return
	}
	if nstate == model.PayFaild {
		// coupon back to user
		if ls, err = s.dao.ConsumeCouponLog(c, o.Mid, o.OrderNo, model.Consume); err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(ls) == 0 {
			err = fmt.Errorf("ConsumeCouponLog not found (mid:%d,orderNo:%s)", o.Mid, o.OrderNo)
			return
		}
		if err = s.UpdateBalance(c, tx, o.Mid, o.CouponType, ls, o.OrderNo); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

// UpdateBalance update user balance.
func (s *Service) UpdateBalance(c context.Context, tx *sql.Tx, mid int64, ct int8, ls []*model.CouponBalanceChangeLog, orderNo string) (err error) {
	var (
		now   = time.Now()
		bs    []*model.CouponBalanceInfo
		aff   int64
		usebs []*model.CouponBalanceInfo
		blogs []*model.CouponBalanceChangeLog
	)
	if bs, err = s.dao.BlanceList(c, mid, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(bs) == 0 {
		err = fmt.Errorf("coupon balance not found (mid:%d ct:%d)", mid, ct)
		return
	}
	for _, ob := range bs {
		for _, l := range ls {
			if ob.BatchToken == l.BatchToken {
				b := new(model.CouponBalanceInfo)
				b.ID = ob.ID
				b.Ver = ob.Ver
				b.Balance = ob.Balance - l.ChangeBalance
				usebs = append(usebs, b)

				blog := new(model.CouponBalanceChangeLog)
				blog.OrderNo = orderNo
				blog.Mid = mid
				blog.BatchToken = ob.BatchToken
				blog.ChangeType = model.ConsumeFaildBack
				blog.Ctime = xtime.Time(now.Unix())
				blog.Balance = b.Balance
				blog.ChangeBalance = -l.ChangeBalance
				blogs = append(blogs, blog)
			}
		}
	}
	if len(ls) != len(usebs) {
		err = fmt.Errorf("coupon balance not found (mid:%d len(ls):%d) len(usebs):%d", mid, len(ls), len(usebs))
		return
	}
	if len(usebs) == 1 {
		b := usebs[0]
		if aff, err = s.dao.UpdateBlance(c, tx, b.ID, mid, b.Ver, b.Balance); err != nil {
			err = errors.WithStack(err)
			return
		}
	} else {
		if aff, err = s.dao.BatchUpdateBlance(c, tx, mid, usebs); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if int(aff) != len(usebs) {
		err = fmt.Errorf("coupon balance back faild mid(%d) order(%s)", mid, orderNo)
		return
	}
	if _, err = s.dao.BatchInsertBlanceLog(c, tx, mid, blogs); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.dao.DelCouponBalancesCache(c, mid, ct)
	return
}

// CheckInUseCoupon check inuse coupon.
func (s *Service) CheckInUseCoupon() {
	var (
		c   = context.TODO()
		cps []*model.CouponInfo
		t   = time.Now().AddDate(0, 0, -1)
		err error
	)
	log.Info("check inuse coupon job start")
	for i := 0; i < 100; i++ {
		if cps, err = s.dao.CouponList(c, int64(i), model.InUse, t); err != nil {
			log.Error("query coupon list(%d,%v) err(%v)", i, t, err)
			return
		}
		log.Info("check inuse coupon job ing size(%d)", len(cps))
		for _, v := range cps {
			var notifyURL string
			if v.CouponType == model.BangumiVideo {
				notifyURL = s.c.Properties.BangumiNotifyURL
			}
			if len(notifyURL) == 0 {
				continue
			}
			// point callback.
			arg := &model.NotifyParam{
				Mid:         v.Mid,
				CouponToken: v.CouponToken,
				NotifyURL:   notifyURL,
				Type:        v.CouponType,
			}
			if err = s.CheckCouponDeliver(c, arg); err != nil {
				log.Error("CheckCouponDeliver fail arg(%v) err(%v)", arg, err)
				continue
			}
		}
		time.Sleep(time.Second * 1)
	}
	log.Info("check inuse coupon job start")
}

// CheckOrderInPayCoupon check order inuse coupon.
func (s *Service) CheckOrderInPayCoupon() {
	var (
		c   = context.TODO()
		cps []*model.CouponOrder
		t   = time.Now().AddDate(0, 0, -1)
		err error
	)
	log.Info("check inuse coupon order job start")
	if cps, err = s.dao.OrderInPay(c, model.InPay, t); err != nil {
		log.Error("query coupon order list(%d,%v) err(%v)", model.InPay, t, err)
		return
	}
	log.Info("check inuse coupon order job ing size(%d)", len(cps))
	for _, v := range cps {
		var notifyURL string
		if v.CouponType == model.Cartoon {
			notifyURL = s.c.Properties.BangumiNotifyURL
		}
		if len(notifyURL) == 0 {
			continue
		}
		// point callback.
		arg := &model.NotifyParam{
			Mid:         v.Mid,
			CouponToken: v.OrderNo,
			NotifyURL:   notifyURL,
			Type:        int64(v.CouponType),
		}
		if err = s.CheckCouponDeliver(c, arg); err != nil {
			log.Error("CheckCouponDeliver order fail arg(%v) err(%v)", arg, err)
			continue
		}
	}
	log.Info("check inuse coupon order job start")
}

// ByOrderNo by order no.
func (s *Service) ByOrderNo(c context.Context, orderNo string) (o *model.CouponOrder, err error) {
	if o, err = s.dao.ByOrderNo(c, orderNo); err != nil {
		err = errors.WithStack(err)
	}
	return
}
