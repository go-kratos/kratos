package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// CartoonCoupon user cartoon coupon.
func (s *Service) CartoonCoupon(c context.Context, mid int64, ct int8) (res []*model.CouponBalanceInfo, err error) {
	var (
		mc  = true
		now = time.Now().Unix()
	)
	if mid <= 0 {
		return
	}
	if res, err = s.dao.CouponBlanceCache(c, mid, ct); err != nil {
		log.Error("mid(%d) err(%+v)", mid, err)
		err = nil
		mc = false
	}
	// empty cache.
	if res != nil && len(res) == 0 {
		return
	}
	tmp := []*model.CouponBalanceInfo{}
	for _, c := range res {
		if c.ExpireTime > now && c.StartTime < now {
			tmp = append(tmp, c)
		}
	}
	if len(tmp) > 0 {
		res = tmp
		return
	}
	if res, err = s.dao.BlanceNoStartCheckList(c, mid, ct, now); err != nil {
		err = errors.WithStack(err)
		return
	} else if len(res) == 0 {
		res = _emptyBlance
	}
	if mc {
		s.dao.SetCouponBlanceCache(c, mid, ct, res)
	}
	tmp = []*model.CouponBalanceInfo{}
	for _, c := range res {
		if c.StartTime < now {
			tmp = append(tmp, c)
		}
	}
	res = tmp
	return

}

// CartoonUse cartoon coupon use.
func (s *Service) CartoonUse(c context.Context, mid int64, thirdTradeNo string, ct int8, useVer int64, remark string, tips string,
	count int64) (ret int8, token string, err error) {
	var (
		o      *model.CouponOrder
		now    = time.Now().Unix()
		blance int64
		cs     []*model.CouponBalanceInfo
		lock   bool
	)
	defer func() {
		if lock {
			s.dao.DelUniqueKey(c, thirdTradeNo, ct)
		}
	}()
	if lock = s.dao.AddUseUniqueLock(c, thirdTradeNo, ct); !lock {
		err = ecode.CouPonUseTooFrequently
		return
	}
	ret = model.UseFaild
	if o, err = s.dao.ByThirdTradeNo(c, thirdTradeNo, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	if o != nil {
		if o.State == model.InPay || o.State == model.PaySuccess {
			ret = model.UseSuccess
		}
		token = o.OrderNo
		return
	}
	if cs, err = s.dao.CouponBlances(c, mid, ct, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(cs) == 0 {
		return
	}
	// check blance
	for _, v := range cs {
		blance += v.Balance
	}
	if blance < count {
		err = ecode.CouPonNotEnoughErr
		return
	}
	var orderNo string
	if orderNo, err = s.ConsumeCoupon(c, mid, ct, cs, count, thirdTradeNo, remark, tips, useVer); err != nil {
		err = errors.WithStack(err)
		return
	}
	token = orderNo
	ret = model.UseSuccess
	return
}

// ConsumeCoupon consume coupon.
func (s *Service) ConsumeCoupon(c context.Context, mid int64, ct int8, cs []*model.CouponBalanceInfo, count int64, thirdTradeNo string,
	remark string, tips string, useVer int64) (orderNo string, err error) {
	var (
		tx  *sql.Tx
		aff int64
		o   *model.CouponOrder
		ol  *model.CouponOrderLog
		now = time.Now()
	)
	orderNo = s.orderID()
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.Wrapf(err, "s.dao.BeginTran(%d)", mid)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				err = errors.Wrapf(err, "tx.Rollback(%d)", mid)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = errors.Wrapf(err, "tx.Commit(%d)", mid)
		}
		s.dao.DelCouponBalancesCache(c, mid, ct)
	}()
	if err = s.UpdateBalance(c, tx, mid, count, cs, orderNo, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	// add order
	o = new(model.CouponOrder)
	o.OrderNo = orderNo
	o.Mid = mid
	o.Count = count
	o.State = model.InPay
	o.CouponType = ct
	o.ThirdTradeNo = thirdTradeNo
	o.Remark = remark
	o.Tips = tips
	o.UseVer = useVer
	o.Ver = 1
	o.Ctime = xtime.Time(now.Unix())
	if aff, err = s.dao.AddOrder(c, tx, o); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = ecode.CouPonConsumeFaildErr
		return
	}
	// add order log.
	ol = new(model.CouponOrderLog)
	ol.OrderNo = orderNo
	ol.Mid = mid
	ol.State = model.InPay
	ol.Ctime = xtime.Time(now.Unix())
	if _, err = s.dao.AddOrderLog(c, tx, ol); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateBalance update user balance.
func (s *Service) UpdateBalance(c context.Context, tx *sql.Tx, mid int64, count int64, cs []*model.CouponBalanceInfo, orderNo string, ct int8) (err error) {
	var (
		now   = time.Now()
		usebs []*model.CouponBalanceInfo
		blogs []*model.CouponBalanceChangeLog
		aff   int64
	)
	for _, v := range cs {
		if v.Balance <= 0 {
			continue
		}
		b := new(model.CouponBalanceInfo)
		b.ID = v.ID
		b.Ver = v.Ver

		blog := new(model.CouponBalanceChangeLog)
		blog.OrderNo = orderNo
		blog.Mid = mid
		blog.BatchToken = v.BatchToken
		blog.ChangeType = model.Consume
		blog.Ctime = xtime.Time(now.Unix())

		if v.Balance >= count {
			b.Balance = v.Balance - count
			usebs = append(usebs, b)
			blog.Balance = b.Balance
			blog.ChangeBalance = -count
			blogs = append(blogs, blog)
			break
		}
		count = count - v.Balance
		b.Balance = 0
		usebs = append(usebs, b)

		blog.Balance = b.Balance
		blog.ChangeBalance = -v.Balance
		blogs = append(blogs, blog)
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
		err = ecode.CouPonConsumeFaildErr
		return
	}
	if _, err = s.dao.BatchInsertBlanceLog(c, tx, mid, blogs); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// CouponCartoonPage coupon cartoon page.
func (s *Service) CouponCartoonPage(c context.Context, mid int64, state int8, pn, ps int) (data *model.CouponCartoonPageResp, err error) {
	var (
		t           = time.Now().Unix()
		bs          []*model.CouponBalanceInfo
		os          []*model.CouponOrder
		res         []*model.CouponPageResp
		stime       = time.Now().AddDate(0, -3, 0)
		couponCount int64
		count       int64
	)
	data = new(model.CouponCartoonPageResp)
	if count, err = s.dao.CouponCarToonCount(c, mid, t, model.CouponCartoon, state, stime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if count <= 0 {
		return
	}
	if ps == 0 {
		ps = _defps
	}
	if pn == 0 {
		pn = _defpn
	}
	switch state {
	case model.NotUsed:
		bs, err = s.dao.CouponNotUsedPage(c, mid, model.CouponCartoon, t, stime, pn, ps)
		res, couponCount = s.convertByBalance(bs, state)
	case model.Used:
		os, err = s.dao.OrderUsedPage(c, mid, model.PaySuccess, model.CouponCartoon, stime, pn, ps)
		res, couponCount = s.convertByOrder(os)
	case model.Expire:
		bs, err = s.dao.CouponExpirePage(c, mid, model.CouponCartoon, t, stime, pn, ps)
		res, couponCount = s.convertByBalance(bs, state)
	default:
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	data.Count = count
	data.CouponCount = couponCount
	data.List = res
	return
}

func (s *Service) convertByBalance(bs []*model.CouponBalanceInfo, state int8) (res []*model.CouponPageResp, couponCount int64) {
	for _, v := range bs {
		r := new(model.CouponPageResp)
		r.ID = v.ID
		switch state {
		case model.NotUsed:
			r.Title = _defCartoonTitle
			r.Time = v.ExpireTime
			r.Count = v.Balance
		case model.Expire:
			r.Title = _defCartoonTitle
			r.Time = v.ExpireTime
			r.Count = v.Balance
		}
		couponCount += v.Balance
		res = append(res, r)
	}
	return
}

func (s *Service) convertByOrder(os []*model.CouponOrder) (res []*model.CouponPageResp, couponCount int64) {
	for _, v := range os {
		r := new(model.CouponPageResp)
		r.ID = v.ID
		r.Title = v.Remark
		r.Tips = v.Tips
		r.Time = v.Mtime.Time().Unix()
		r.Count = v.Count
		couponCount += v.Count
		res = append(res, r)
	}
	return
}

// AddCartoonCoupon add cartoon coupon.
func (s *Service) AddCartoonCoupon(c context.Context, bi *model.CouponBatchInfo, mid int64, ct int64, origin int64, count int) (err error) {
	var (
		tx   *sql.Tx
		oldb *model.CouponBalanceInfo
		logs []*model.CouponBalanceChangeLog
		hc   int64
	)
	// limit count
	if bi.LimitCount != _maxCount {
		// check grant state
		if logs, err = s.dao.GrantCouponLog(c, mid, bi.BatchToken, model.VipSalary); err != nil {
			err = errors.WithStack(err)
			return
		}
		for _, v := range logs {
			hc += v.ChangeBalance
		}
		if int64(count)+hc > bi.LimitCount {
			err = ecode.CouPonBatchLimitErr
			return
		}
	}
	if oldb, err = s.dao.ByMidAndBatchToken(c, mid, bi.BatchToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.Wrapf(err, "s.dao.BeginTran(%d)", mid)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				err = errors.Wrapf(err, "tx.Rollback(%d)", mid)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = errors.Wrapf(err, "tx.Commit(%d)", mid)
		}
		s.dao.DelCouponBalancesCache(c, mid, int8(ct))
	}()
	if err = s.UpdateBranch(c, tx, count, bi); err != nil {
		err = errors.WithStack(err)
		return
	}
	bc := new(model.CouponBalanceInfo)
	bc.BatchToken = bi.BatchToken
	bc.Mid = mid
	bc.Balance = int64(count)
	bc.StartTime = bi.StartTime
	bc.ExpireTime = bi.ExpireTime
	bc.Origin = origin
	bc.CouponType = ct
	bc.Ver = 1
	bc.CTime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddBalanceCoupon(c, tx, bc); err != nil {
		err = errors.WithStack(err)
		return
	}
	blog := new(model.CouponBalanceChangeLog)
	blog.Mid = mid
	blog.BatchToken = bi.BatchToken
	if oldb == nil {
		blog.Balance = int64(count)
	} else {
		blog.Balance = oldb.Balance + int64(count)
	}
	blog.ChangeBalance = int64(count)
	blog.ChangeType = model.VipSalary
	blog.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddBalanceChangeLog(c, tx, blog); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// orderID get order id.
func (s *Service) orderID() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%010d", s.r.Int63n(9999999999)))
	b.WriteString(time.Now().Format("150405"))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	return b.String()
}

// CarToonCouponCount cartoon coupon count.
func (s *Service) CarToonCouponCount(c context.Context, mid int64, ct int8) (count int, err error) {
	var res []*model.CouponBalanceInfo
	if res, err = s.CartoonCoupon(c, mid, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range res {
		count += int(v.Balance)
	}
	return
}
