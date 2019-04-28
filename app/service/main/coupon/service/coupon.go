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

// UserCoupon user coupon.
func (s *Service) UserCoupon(c context.Context, mid int64, ct int8) (cs []*model.CouponInfo, err error) {
	var (
		mc  = true
		now = time.Now().Unix()
	)
	if mid <= 0 {
		return
	}
	if cs, err = s.dao.CouponsCache(c, mid, ct); err != nil {
		log.Error("mid(%d) err(%+v)", mid, err)
		err = nil
		mc = false
	}
	// empty cache.
	if cs != nil && len(cs) == 0 {
		return
	}
	tmp := []*model.CouponInfo{}
	for _, c := range cs {
		if c.ExpireTime > now && c.StartTime < now {
			tmp = append(tmp, c)
		}
	}
	if len(tmp) > 0 {
		cs = tmp
		return
	}
	if cs, err = s.dao.CouponNoStartCheckList(c, mid, model.NotUsed, ct, now); err != nil {
		err = errors.WithStack(err)
		return
	} else if len(cs) == 0 {
		cs = _emptyCoupons
	}
	if mc {
		s.dao.SetCouponsCache(c, mid, ct, cs)
	}
	tmp = []*model.CouponInfo{}
	for _, c := range cs {
		if c.StartTime < now {
			tmp = append(tmp, c)
		}
	}
	cs = tmp
	return
}

// UseCoupon use coupon.
func (s *Service) UseCoupon(c context.Context, mid int64, oid int64, remark string, orderNO string, ct int8, useVer int64) (ret int8, token string, err error) {
	var (
		cp   *model.CouponInfo
		now  = time.Now().Unix()
		cs   []*model.CouponInfo
		lock bool
	)
	defer func() {
		if lock {
			s.dao.DelUniqueKey(c, orderNO, ct)
		}
	}()
	if lock = s.dao.AddUseUniqueLock(c, orderNO, ct); !lock {
		err = ecode.CouPonUseTooFrequently
		return
	}
	ret = model.UseFaild
	if cp, err = s.dao.ByOrderNO(c, mid, orderNO, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp != nil {
		if cp.State == model.InUse || cp.State == model.Used {
			ret = model.UseSuccess
		}
		token = cp.CouponToken
		return
	}
	if cs, err = s.dao.CouponList(c, mid, model.NotUsed, ct, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(cs) == 0 {
		return
	}
	cp = cs[0]
	if err = s.UpdateCoupon(c, cp, oid, remark, orderNO, ct, useVer); err != nil {
		err = errors.Wrapf(err, "use coupon error(%d)", mid)
		return
	}
	token = cp.CouponToken
	ret = model.UseSuccess
	return
}

// UpdateCoupon update coupon info.
func (s *Service) UpdateCoupon(c context.Context, cp *model.CouponInfo, oid int64, remark string, orderNO string, ct int8, useVer int64) (err error) {
	var (
		tx  *sql.Tx
		aff int64
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("%+v", err)
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
		s.dao.DelCouponsCache(c, cp.Mid, ct)
	}()
	cp.State = model.InUse
	cp.OrderNO = orderNO
	cp.Oid = oid
	cp.Remark = remark
	cp.UseVer = useVer
	if aff, err = s.dao.UpdateCouponInUse(c, tx, cp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("use coupon faild")
		return
	}
	l := &model.CouponChangeLog{}
	l.CouponToken = cp.CouponToken
	l.Mid = cp.Mid
	l.State = int8(cp.State)
	l.Ctime = xtime.Time(time.Now().Unix())
	if aff, err = s.dao.InsertPointHistory(c, tx, l); err != nil {
		log.Error("%+v", err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("add change log faild")
		return
	}
	return
}

// CouponInfo coupon info.
func (s *Service) CouponInfo(c context.Context, mid int64, token string) (cp *model.CouponInfo, err error) {
	if cp, err = s.dao.CouponInfo(c, mid, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp == nil {
		err = ecode.NothingFound
		return
	}
	if cp.State == model.NotUsed && cp.ExpireTime < time.Now().Unix() {
		cp.State = model.Expire
	}
	return
}

// CouponPage coupon page.
func (s *Service) CouponPage(c context.Context, mid int64, state int8, pn, ps int) (count int64, res []*model.CouponPageResp, err error) {
	var (
		t  = time.Now().Unix()
		cp []*model.CouponInfo
	)
	stime := time.Now().AddDate(0, -3, 0)
	if count, err = s.dao.CountByState(c, mid, state, t, stime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if count == 0 {
		return
	}
	if pn == 0 {
		pn = _defpn
	}
	if ps == 0 {
		ps = _defps
	}
	if cp, err = s.dao.CouponPage(c, mid, state, t, (pn-1)*ps, ps, stime); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(cp) == 0 {
		return
	}
	for _, v := range cp {
		r := &model.CouponPageResp{}
		r.ID = v.ID
		switch state {
		case model.NotUsed:
			r.Title = _deftitle
			r.Time = v.ExpireTime
		case model.Used:
			r.Title = v.Remark
			r.Time = v.MTime.Time().Unix()
		case model.Expire:
			r.Title = _deftitle
			r.Time = v.ExpireTime
		}
		r.RefID = v.Oid
		res = append(res, r)
	}
	return
}

// AddCoupon add coupon.
func (s *Service) AddCoupon(c context.Context, mid int64, startTime int64, expireTime int64, ct int64, origin int64) (err error) {
	cp := &model.CouponInfo{}
	cp.CouponToken = s.token()
	cp.Mid = mid
	cp.State = model.NotUsed
	cp.StartTime = startTime
	cp.ExpireTime = expireTime
	cp.Origin = origin
	cp.CouponType = ct
	cp.CTime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddCoupon(c, cp); err != nil {
		log.Error("%+v", err)
		return
	}
	s.dao.DelCouponsCache(c, mid, int8(ct))
	return
}

// get coupon token
func (s *Service) token() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%07d", s.r.Int63n(9999999)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("20060102150405"))
	return b.String()
}

//ChangeState change state.
func (s *Service) ChangeState(c context.Context, mid int64, userVer int64, ver int64, couponToken string) (err error) {
	if _, err = s.dao.UpdateCoupon(c, mid, model.Used, userVer, ver, couponToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// SalaryCoupon salary coupon.
func (s *Service) SalaryCoupon(c context.Context, a *model.ArgSalaryCoupon) (err error) {
	var (
		bi   *model.CouponBatchInfo
		lock bool
		htc  int
	)
	defer func() {
		if lock {
			s.dao.DelGrantKey(c, a.BatchToken, a.Mid)
		}
	}()
	if a.Count <= 0 || a.Count > model.MaxSalaryCount {
		err = ecode.RequestErr
		return
	}
	if bi = s.allBranchInfo[a.BatchToken]; bi == nil {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	if bi.State == model.BatchStateBlock {
		err = ecode.CouponBatchBlockErr
		return
	}
	// max count
	if bi.MaxCount != _maxCount {
		if htc, err = s.CurrentCount(c, a.BatchToken); err != nil {
			err = errors.WithStack(err)
			return
		}
		if int64(htc+a.Count) > bi.MaxCount {
			err = ecode.CouPonBatchNotEnoughErr
			return
		}
	}
	// limit count
	if bi.LimitCount != _maxCount {
		if lock = s.dao.AddGrantUniqueLock(c, a.BatchToken, a.Mid); !lock {
			err = ecode.CouPonGrantTooFrequently
			return
		}
	}
	switch a.CouponType {
	case model.CouponVideo:
		if err = s.AddVideoCoupon(c, bi, a.Mid, a.CouponType, a.Origin, a.Count); err != nil {
			err = errors.WithStack(err)
			return
		}
	case model.CouponCartoon:
		if err = s.AddCartoonCoupon(c, bi, a.Mid, a.CouponType, a.Origin, a.Count); err != nil {
			err = errors.WithStack(err)
			return
		}
	case model.CouponAllowance:
		if err = s.AddAllowanceCoupon(c, bi, a.Mid, a.Count, a.Origin, a.AppID); err != nil {
			err = errors.WithStack(err)
			return
		}
	default:
		err = ecode.CouPonTypeNotExistErr
	}
	return
}

// SalaryCouponForThird salary coupon for third.
func (s *Service) SalaryCouponForThird(c context.Context, a *model.ArgSalaryCoupon) (res *model.SalaryCouponForThirdResp, err error) {
	var lock bool
	if lock = s.dao.AddUniqueNoLock(c, a.UniqueNo); !lock {
		err = ecode.VipCouponUniqueNoErr
		return
	}
	if err = s.SalaryCoupon(c, a); err != nil {
		return
	}
	var bi *model.CouponBatchInfo
	if bi = s.allBranchInfo[a.BatchToken]; bi == nil {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	res = new(model.SalaryCouponForThirdResp)
	res.Amount = bi.Amount
	res.FullAmount = bi.FullAmount
	res.Description = bi.Name
	return
}

// AddVideoCoupon add video coupon.
func (s *Service) AddVideoCoupon(c context.Context, bi *model.CouponBatchInfo, mid int64, ct int64, origin int64, count int) (err error) {
	var (
		tx *sql.Tx
		hc int64
	)
	// limit count
	if bi.LimitCount != _maxCount {
		// check grant state
		if hc, err = s.dao.CountByBranchToken(c, mid, bi.BatchToken); err != nil {
			err = errors.WithStack(err)
			return
		}
		if int64(count)+hc > bi.LimitCount {
			err = ecode.CouPonBatchLimitErr
			return
		}
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
		s.dao.DelCouponsCache(c, mid, int8(ct))
	}()
	if err = s.UpdateBranch(c, tx, count, bi); err != nil {
		err = errors.WithStack(err)
		return
	}
	cps := make([]*model.CouponInfo, count)
	for i := 0; i < count; i++ {
		cp := &model.CouponInfo{}
		cp.CouponToken = s.token()
		cp.Mid = mid
		cp.State = model.NotUsed
		cp.StartTime = bi.StartTime
		cp.ExpireTime = bi.ExpireTime
		cp.Origin = origin
		cp.CouponType = ct
		cp.CTime = xtime.Time(time.Now().Unix())
		cp.BatchToken = bi.BatchToken
		cps[i] = cp
	}
	if _, err = s.dao.BatchAddCoupon(c, tx, mid, cps); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// VideoCouponCount video coupon count.
func (s *Service) VideoCouponCount(c context.Context, mid int64, ct int8) (count int, err error) {
	var cs []*model.CouponInfo
	if cs, err = s.UserCoupon(c, mid, ct); err != nil {
		err = errors.WithStack(err)
		return
	}
	count = len(cs)
	return
}
