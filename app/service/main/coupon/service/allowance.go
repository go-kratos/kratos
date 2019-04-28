package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	coupinv1 "go-common/app/service/main/coupon/api"
	"go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vipinfo/api"
	vimdl "go-common/app/service/main/vipinfo/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// AllowanceCoupon allowance coupon.
func (s *Service) AllowanceCoupon(c context.Context, a *model.ArgAllowanceCoupons) (res []*model.CouponAllowanceInfo, err error) {
	var (
		mc  = true
		now = time.Now().Unix()
	)
	if a.Mid <= 0 {
		return
	}
	if res, err = s.dao.CouponAllowanceCache(c, a.Mid, a.State); err != nil {
		log.Error("mid(%d) err(%+v)", a.Mid, err)
		err = nil
		mc = false
	}
	if res == nil {
		// not found
		if res, err = s.dao.ByStateAndExpireAllowances(c, a.Mid, a.State, now); err != nil {
			err = errors.WithStack(err)
			return
		} else if len(res) == 0 {
			res = _emptyAllowance
		}
		if mc {
			s.dao.SetCouponAllowanceCache(c, a.Mid, a.State, res)
		}
	}
	tmp := []*model.CouponAllowanceInfo{}
	for _, c := range res {
		if c.ExpireTime > now {
			tmp = append(tmp, c)
		}
	}
	res = tmp
	return
}

// JudgeCouponUsable judge coupon is able.
func (s *Service) JudgeCouponUsable(c context.Context, mid int64, price float64, couponToken string, plat int, prodLimMonth, prodLimRenewal int8) (cp *model.CouponAllowanceInfo, err error) {
	var (
		bi  *model.CouponBatchInfo
		now = time.Now().Unix()
	)
	if cp, err = s.dao.AllowanceByToken(c, mid, couponToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp == nil {
		err = ecode.CouPonTokenNotFoundErr
		return
	}
	if bi = s.allBranchInfo[cp.BatchToken]; bi == nil {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	if bi.State != model.BatchStateNormal {
		err = ecode.CouPonHadBlockErr
		return
	}
	if cp.State != model.NotUsed {
		err = ecode.CouPonHadUseErr
		return
	}
	if cp.StartTime > now || now > cp.ExpireTime {
		err = ecode.CouPonHadExpireErr
		return
	}
	if cp.FullAmount > price {
		err = ecode.CouPonNotFullPriceErr
		return
	}
	if !s.platformLimit(bi.PlatformLimit, plat) {
		err = ecode.CouPonPlatformNotSupportErr
		return
	}
	if !s.productLimit(bi.ProdLimMonth, bi.ProdLimRenewal, prodLimMonth, prodLimRenewal) {
		log.Error("s.productLimit(mid:%d bi:%+v %d %d)", mid, bi, prodLimMonth, prodLimRenewal)
		err = ecode.CouPonProductNotSupportErr
		return
	}
	return
}

// UseAllowanceCoupon use allance coupon.
func (s *Service) UseAllowanceCoupon(c context.Context, arg *model.ArgUseAllowance) (err error) {
	var (
		cp    *model.CouponAllowanceInfo
		exist *model.CouponAllowanceInfo
	)
	if cp, err = s.JudgeCouponUsable(c, arg.Mid, arg.Price, arg.CouponToken, model.PlatformByName[arg.Platform], arg.ProdLimMonth, arg.ProdLimRenewal); err != nil {
		return
	}
	if exist, err = s.dao.AllowanceByOrderNO(c, arg.Mid, arg.OrderNO); err != nil {
		err = errors.WithStack(err)
		return
	}
	if exist != nil {
		err = ecode.CouPonOrderHadUseErr
		return
	}
	if err = s.UpdateAllowanceCoupon(c, cp, model.InUse, arg.Remark, arg.OrderNO, model.AllowanceConsume); err != nil {
		err = errors.Wrapf(err, "use allowance coupon error(%d)", arg.Mid)
		return
	}
	return
}

// UpdateAllowanceCoupon update coupon info.
func (s *Service) UpdateAllowanceCoupon(c context.Context, cp *model.CouponAllowanceInfo, state int32, remark string, orderNO string, changeType int8) (err error) {
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
				log.Error("tx.Rollback %+v", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit %+v", err)
		}
		s.dao.DelCouponAllowancesKey(c, cp.Mid, model.NotUsed)
		s.dao.DelCouponAllowancesKey(c, cp.Mid, model.InUse)
	}()
	cp.State = state
	cp.OrderNO = orderNO
	cp.Remark = remark
	if aff, err = s.dao.UpdateAllowanceCouponInUse(c, tx, cp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("use coupon faild")
		return
	}
	l := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		Mid:         cp.Mid,
		State:       int8(cp.State),
		Ctime:       xtime.Time(time.Now().Unix()),
		OrderNO:     orderNO,
		ChangeType:  changeType,
	}
	if aff, err = s.dao.InsertCouponAllowanceHistory(c, tx, l); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("add change log faild")
		return
	}
	return
}

// UsableAllowanceCoupon usable coupon.
func (s *Service) UsableAllowanceCoupon(c context.Context, mid int64, price float64, plat int, prodLimMonth, prodLimRenewal int8) (res *model.CouponAllowancePanelInfo, err error) {
	var (
		us  []*model.CouponAllowancePanelInfo
		all []*model.CouponAllowanceInfo
	)
	if all, err = s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   mid,
		State: model.NotUsed,
	}); err != nil {
		return
	}
	if us, _, _, err = s.UsableAllowanceCoupons(c, mid, price, all, plat, prodLimMonth, prodLimRenewal); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(us) > 0 {
		res = us[0]
	}
	return
}

// MultiUsableAllowanceCoupon mulyi usable allowance coupon.
func (s *Service) MultiUsableAllowanceCoupon(c context.Context, mid int64, price []float64, plat int, prodLimMonth, prodLimRenewal int8) (res map[float64]*model.CouponAllowancePanelInfo, err error) {
	var (
		all []*model.CouponAllowanceInfo
		us  []*model.CouponAllowancePanelInfo
	)
	if len(price) == 0 {
		return
	}
	if all, err = s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   mid,
		State: model.NotUsed,
	}); err != nil {
		return
	}
	res = make(map[float64]*model.CouponAllowancePanelInfo, len(price))
	for _, v := range price {
		if us, _, _, err = s.UsableAllowanceCoupons(c, mid, v, all, plat, prodLimMonth, prodLimRenewal); err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(us) == 0 {
			continue
		}
		res[v] = us[0]
	}
	return
}

// UsableAllowanceCoupons usable allowance coupon list.
func (s *Service) UsableAllowanceCoupons(c context.Context, mid int64, price float64, all []*model.CouponAllowanceInfo, plat int, prodLimMonth, prodLimRenewal int8) (us []*model.CouponAllowancePanelInfo,
	ds []*model.CouponAllowancePanelInfo, ui []*model.CouponAllowancePanelInfo, err error) {
	var (
		bi  *model.CouponBatchInfo
		now = time.Now().Unix()
	)
	us = make([]*model.CouponAllowancePanelInfo, 0)
	ds = make([]*model.CouponAllowancePanelInfo, 0)
	ui = make([]*model.CouponAllowancePanelInfo, 0)
	if len(all) == 0 {
		return
	}
	for _, r := range all {
		if bi = s.allBranchInfo[r.BatchToken]; bi == nil {
			continue
		}
		if bi.State == model.BatchStateBlock || bi.State == model.Block {
			continue
		}
		ok := true
		explains := []string{}
		if !s.platformLimit(bi.PlatformLimit, plat) {
			explains = append(explains, model.CouponPlatformExplain)
			ok = false
		}
		if !s.productLimit(bi.ProdLimMonth, bi.ProdLimRenewal, prodLimMonth, prodLimRenewal) {
			log.Warn("s.productLimit(%+v %d %d)", bi, prodLimMonth, prodLimRenewal)
			explains = append(explains, model.CouponProductExplain)
			ok = false
		}
		if r.FullAmount > price {
			explains = append(explains, model.CouponFullAmountDissatisfy)
			ok = false
		}
		if r.StartTime > now {
			explains = append(explains, model.CouponNotInUsableTime)
			ok = false
		}
		if r.State == model.InUse {
			if len(explains) > 0 {
				// can not un block
				r.State = model.NotUsed
			}
			if ok {
				ui = append(ui, s.convertCoupon(r, explains, price, model.AllowanceDisables))
			}
			ok = false
		}
		if ok {
			us = append(us, s.convertCoupon(r, explains, price, model.AllowanceUsable))
		} else {
			ds = append(ds, s.convertCoupon(r, explains, price, model.AllowanceDisables))
		}
	}
	if len(us) > 0 {
		sort.Slice(us, func(i int, j int) bool {
			return us[i].ExpireTime < us[j].ExpireTime
		})
		sort.Slice(us, func(i int, j int) bool {
			return us[i].Amount > us[j].Amount
		})
		us[0].Selected = model.Seleted
	}
	return
}

// AllowancePanelCoupons usable allowance coupon list.
func (s *Service) AllowancePanelCoupons(c context.Context, mid int64, price float64, plat int, prodLimMonth, prodLimRenewal int8) (us []*model.CouponAllowancePanelInfo,
	ds []*model.CouponAllowancePanelInfo, ui []*model.CouponAllowancePanelInfo, err error) {
	var (
		now   = time.Now().Unix()
		all   []*model.CouponAllowanceInfo
		inuse []*model.CouponAllowanceInfo
	)
	if all, err = s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   mid,
		State: model.NotUsed,
	}); err != nil {
		return
	}
	// coupon using
	if inuse, err = s.dao.ByStateAndExpireAllowances(c, mid, model.InUse, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	all = append(all, inuse...)
	if us, ds, ui, err = s.UsableAllowanceCoupons(c, mid, price, all, plat, prodLimMonth, prodLimRenewal); err != nil {
		err = errors.WithStack(err)
		return
	}
	sort.Slice(ds, func(i int, j int) bool {
		if ds[i].State == ds[j].State {
			return ds[i].Amount > ds[j].Amount
		}
		return ds[i].State > ds[j].State
	})
	return
}

func (s *Service) convertCoupon(c *model.CouponAllowanceInfo, explains []string, price float64, usable int8) (r *model.CouponAllowancePanelInfo) {
	var (
		bi *model.CouponBatchInfo
	)
	r = new(model.CouponAllowancePanelInfo)
	r.CouponToken = c.CouponToken
	r.Amount = c.Amount
	r.State = c.State
	r.FullLimitExplain = fmt.Sprintf(model.CouponFullAmountLimit, strconv.FormatFloat(c.FullAmount, 'f', 0, 64))
	r.FullAmount = c.FullAmount
	r.StartTime = c.StartTime
	r.ExpireTime = c.ExpireTime
	if bi = s.allBranchInfo[c.BatchToken]; bi != nil {
		r.ScopeExplainFmt(bi.PlatformLimit, bi.ProdLimMonth, bi.ProdLimRenewal, s.c.Platform)
	}
	r.CouponDiscountPrice = price - c.Amount
	r.DisablesExplains = strings.Join(explains, ",")
	r.OrderNO = c.OrderNO
	r.Name = model.CouponAllowanceName
	r.Usable = usable
	return
}

// AddAllowanceCoupon add allowance coupon.
func (s *Service) AddAllowanceCoupon(c context.Context, bi *model.CouponBatchInfo, mid int64, count int, origin int64, appID int64) (err error) {
	var (
		tx *sql.Tx
		hc int64
	)
	// limit count
	if bi.LimitCount != _maxCount {
		// check grant state
		if hc, err = s.dao.CountByAllowanceBranchToken(c, mid, bi.BatchToken); err != nil {
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
		s.dao.DelCouponAllowancesKey(c, mid, model.NotUsed)
	}()
	if err = s.UpdateBranch(c, tx, count, bi); err != nil {
		err = errors.WithStack(err)
		return
	}
	cps := make([]*model.CouponAllowanceInfo, count)
	stime := bi.StartTime
	etime := bi.ExpireTime
	if bi.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		stime = now.Unix()
		etime = now.AddDate(0, 0, int(bi.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		if bi.ExpireTime < time.Now().Unix() {
			err = ecode.CouponBatchExpireTimeErr
			return
		}
	}
	for i := 0; i < count; i++ {
		cps[i] = &model.CouponAllowanceInfo{
			CouponToken: s.token(),
			Mid:         mid,
			State:       model.NotUsed,
			StartTime:   stime,
			ExpireTime:  etime,
			Origin:      origin,
			CTime:       xtime.Time(time.Now().Unix()),
			BatchToken:  bi.BatchToken,
			Amount:      bi.Amount,
			FullAmount:  bi.FullAmount,
			AppID:       appID,
		}
	}
	if _, err = s.dao.BatchAddAllowanceCoupon(c, tx, mid, cps); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//AllowanceInfo get allowance info by token.
func (s *Service) AllowanceInfo(c context.Context, mid int64, token string) (res *model.CouponAllowanceInfo, err error) {
	if res, err = s.dao.AllowanceByToken(c, mid, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		err = ecode.CouPonTokenNotFoundErr
		return
	}
	return
}

//CancelUseCoupon cancel use coupon.
func (s *Service) CancelUseCoupon(c context.Context, mid int64, token string) (err error) {
	var (
		cp *model.CouponAllowanceInfo
	)
	if cp, err = s.AllowanceInfo(c, mid, token); err != nil {
		return
	}
	if cp == nil {
		err = ecode.CouPonTokenNotFoundErr
		return
	}
	if cp.State != model.InUse {
		err = ecode.CouPonStateCanNotCancelErr
		return
	}
	if err = s.UpdateAllowanceCoupon(c, cp, model.NotUsed, "", "", model.AllowanceCancel); err != nil {
		err = errors.Wrapf(err, "cancel allowance coupon error(%d)", mid)
		return
	}
	return
}

//CouponNotify coupon notify.
func (s *Service) CouponNotify(c context.Context, mid int64, orderNo string, payState int8) (err error) {
	var (
		state      int32
		changeType int8
		cp         *model.CouponAllowanceInfo
		remark     string
	)
	if cp, err = s.dao.AllowanceByOrderNO(c, mid, orderNo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp == nil {
		err = ecode.CouPonTokenNotFoundErr
		return
	}
	if cp.State != model.InUse {
		err = ecode.CouPonStateCanNotCancelErr
		return
	}
	switch payState {
	case model.AllowanceUseFaild:
		state = model.NotUsed
		changeType = model.AllowanceConsumeFaild
		orderNo = ""
		remark = ""
	case model.AllowanceUseSuccess:
		state = model.Used
		changeType = model.AllowanceConsumeSuccess
		orderNo = cp.OrderNO
		remark = cp.Remark
	default:
		err = ecode.CouPonNotifyStateErr
		return
	}
	if err = s.UpdateAllowanceCoupon(c, cp, state, remark, orderNo, changeType); err != nil {
		err = errors.Wrapf(err, "notify allowance coupon error(%d)", mid)
		return
	}
	return
}

// AllowanceList allowance list.
func (s *Service) AllowanceList(c context.Context, mid int64, state int8) (res []*model.CouponAllowancePanelInfo, err error) {
	var (
		t     = time.Now().Unix()
		stime = time.Now().AddDate(0, -3, 0)
		list  []*model.CouponAllowanceInfo
		es    []string
	)
	if list, err = s.dao.AllowanceList(c, mid, state, t, stime); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range list {
		if state == model.Expire {
			v.State = model.Expire
		}
		res = append(res, s.convertCoupon(v, es, float64(0), model.AllowanceDisables))
	}
	return
}

func (s *Service) platformLimit(pstr string, plat int) (ok bool) {
	var (
		ps  []string
		err error
		p   int
	)
	if len(pstr) == 0 {
		ok = true
		return
	}
	// if mapping, success := model.PlatfromMapping[plat]; success {
	// 	plat = mapping
	// }
	ps = strings.Split(pstr, ",")
	for _, v := range ps {
		if p, err = strconv.Atoi(v); err != nil {
			continue
		}
		if plat == p {
			ok = true
			break
		}
	}
	return
}

// productLimit 商品限制验证.
func (s *Service) productLimit(bplm, bplr, plm, plr int8) (ok bool) {
	if bplm == model.None && bplr == model.None {
		return true
	}
	if bplm == model.None && bplr == plr {
		return true
	}
	if bplr == model.None && bplm == plm {
		return true
	}
	if bplm == plm && bplr == plr {
		return true
	}
	return
}

// UseNotify 同步检查coupon是否可用.
func (s *Service) UseNotify(c context.Context, arg *model.ArgAllowanceCheck) (cp *model.CouponAllowanceInfo, err error) {
	// 1.检查order是否绑定了券
	if cp, err = s.dao.GetCouponByOrderNo(c, arg.Mid, arg.OrderNo); err != nil {
		return
	}
	if cp == nil || cp.Mid != arg.Mid {
		err = ecode.CouPonUseFaildErr
		return
	}
	if cp.State == model.Used {
		return
	}
	// 2.检查券是否可标记为已使用
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
			if theErr := tx.Rollback(); theErr != nil {
				log.Error("tx.Rollback %+v", theErr)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit %+v", err)
		}
		s.dao.DelCouponAllowancesKey(c, arg.Mid, model.NotUsed)
		s.dao.DelCouponAllowancesKey(c, arg.Mid, model.InUse)
	}()
	log.Info("s.dao.UpdateAllowanceCouponToUse(%+v)", cp)
	cp.State = model.Used
	if aff, err = s.dao.UpdateAllowanceCouponToUse(c, tx, cp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		log.Info("s.dao.UpdateAllowanceCouponToUse aff(%s, %d)", arg.OrderNo, aff)

		err = ecode.CouPonHadUseErr
		return
	}
	l := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		Mid:         cp.Mid,
		State:       int8(cp.State),
		Ctime:       xtime.Time(time.Now().Unix()),
		OrderNO:     arg.OrderNo,
		ChangeType:  model.AllowanceConsumeSuccess,
	}
	if aff, err = s.dao.InsertCouponAllowanceHistory(c, tx, l); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("add change log faild")
		return
	}
	return
}

//ReceiveAllowance receive allowance.
func (s *Service) ReceiveAllowance(c context.Context, arg *model.ArgReceiveAllowance) (couponToken string, err error) {
	var (
		bi      *model.CouponBatchInfo
		rlog    *model.CouponReceiveLog
		count   int64
		success bool
	)
	if rlog, err = s.dao.ReceiveLog(c, arg.Appkey, arg.OrderNo, model.CouponAllowance); err != nil {
		err = errors.WithStack(err)
		return
	}
	if rlog != nil {
		log.Info("receive allowance already handler arg:%+v rlog:%+v", arg, rlog)
		couponToken = rlog.CouponToken
		return
	}
	if succeed := s.dao.AddReceiveUniqueLock(c, arg.Appkey, arg.OrderNo, model.CouponAllowance); !succeed {
		log.Info("receive allowance handlering...............arg:%+v", arg)
		return
	}
	defer func() {
		s.dao.DelReceiveUniqueLock(c, arg.Appkey, arg.OrderNo, model.CouponAllowance)
	}()
	if bi, success = s.allBranchInfo[arg.BatchToken]; !success {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	if bi.State == model.BatchStateBlock {
		err = ecode.CouponBatchBlockErr
		return
	}
	if bi.MaxCount >= 0 && bi.MaxCount < bi.CurrentCount+1 {
		err = ecode.CouPonGrantMaxCountErr
		return
	}
	if bi.LimitCount >= 0 {
		if count, err = s.dao.CountByAllowanceBranchToken(c, arg.Mid, bi.BatchToken); err != nil {
			err = errors.WithStack(err)
			return
		}
		if bi.LimitCount >= 0 && count+1 > bi.LimitCount {
			err = ecode.CouPonBatchLimitErr
			return
		}
	}
	if couponToken, err = s.receiveAllowance(c, arg, bi); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.sendMessage(func() {
		s.dao.SendMessage(context.TODO(), strconv.FormatInt(arg.Mid, 10), model.ReceiveMessage, model.ReceiveMessageTitle)
	})
	return
}

func (s *Service) receiveAllowance(c context.Context, arg *model.ArgReceiveAllowance, bi *model.CouponBatchInfo) (couponToken string, err error) {
	var (
		tx  *sql.Tx
		eff int64
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
		s.dao.DelCouponAllowancesKey(c, arg.Mid, model.NotUsed)
	}()
	info := new(model.CouponAllowanceInfo)
	info.BatchToken = bi.BatchToken
	info.Mid = arg.Mid
	info.CouponToken = s.token()
	info.State = model.NotUsed
	info.FullAmount = bi.FullAmount
	info.AppID = bi.AppID
	info.Amount = bi.Amount
	info.Origin = model.AllowanceBusinessReceive
	if bi.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		info.StartTime = time.Now().Unix()
		info.ExpireTime = now.AddDate(0, 0, int(bi.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		if bi.ExpireTime < time.Now().Unix() {
			err = ecode.CouponBatchExpireTimeErr
			return
		}
		info.StartTime = bi.StartTime
		info.ExpireTime = bi.ExpireTime

	}
	if err = s.dao.TxAddAllowanceCoupon(tx, info); err != nil {
		err = errors.WithStack(err)
		return
	}
	changeLog := new(model.CouponAllowanceChangeLog)
	changeLog.State = model.NotUsed
	changeLog.Mid = arg.Mid
	changeLog.CouponToken = info.CouponToken
	changeLog.ChangeType = model.AllowanceReceive
	changeLog.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.InsertCouponAllowanceHistory(c, tx, changeLog); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = s.dao.UpdateBatchInfo(c, tx, bi.BatchToken, 1); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff <= 0 {
		err = ecode.CouponReceiveErr
		return
	}
	rlog := new(model.CouponReceiveLog)
	rlog.Mid = arg.Mid
	rlog.OrderNo = arg.OrderNo
	rlog.Appkey = arg.Appkey
	rlog.CouponToken = info.CouponToken
	rlog.CouponType = model.CouponAllowance
	if err = s.dao.TxAddReceiveLog(tx, rlog); err != nil {
		err = errors.WithStack(err)
		return
	}
	couponToken = info.CouponToken
	return
}

// addPrizeCoupon add allowance coupon.
func (s *Service) addPrizeCoupon(c context.Context, bi *model.CouponBatchInfo, mid, actID int64, cardType int8) (res *model.CouponUserCard, err error) {
	var (
		tx *sql.Tx
	)
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
		s.dao.DelCouponAllowancesKey(c, mid, model.NotUsed)
		s.dao.DelPrizeCardsKey(c, mid, actID)
	}()
	if err = s.UpdateBranch(c, tx, 1, bi); err != nil {
		err = errors.WithStack(err)
		return
	}
	cps := make([]*model.CouponAllowanceInfo, 0)
	stime := bi.StartTime
	etime := bi.ExpireTime
	if bi.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		stime = now.Unix()
		etime = now.AddDate(0, 0, int(bi.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		if bi.ExpireTime < time.Now().Unix() {
			err = ecode.CouponBatchExpireTimeErr
			return
		}
	}
	couponToken := s.token()
	cps = append(cps, &model.CouponAllowanceInfo{
		CouponToken: couponToken,
		Mid:         mid,
		State:       model.NotUsed,
		StartTime:   stime,
		ExpireTime:  etime,
		Origin:      model.AllowanceBusinessNewYear,
		CTime:       xtime.Time(time.Now().Unix()),
		BatchToken:  bi.BatchToken,
		Amount:      bi.Amount,
		FullAmount:  bi.FullAmount,
		AppID:       1,
	})
	var aff int64
	if aff, err = s.dao.BatchAddAllowanceCoupon(c, tx, mid, cps); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = ecode.CouponNewYearGrantErr
		return
	}
	res = &model.CouponUserCard{
		MID:         mid,
		CardType:    cardType,
		State:       model.CardStateOpened,
		BatchToken:  bi.BatchToken,
		CouponToken: couponToken,
		ActID:       actID,
	}
	if aff, err = s.dao.AddUserCard(c, tx, res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = ecode.CouponNewYearGrantErr
		return
	}
	return
}

// PrizeCards .
func (s *Service) PrizeCards(c context.Context, mid int64) (res []*model.PrizeCardRep, err error) {
	var now = time.Now().Unix()
	if s.c.NewYearConf.StartTime > now {
		err = ecode.CouponNewYearNotStartErr
		return
	}
	res, err = s.dao.PrizeCardsCache(c, mid, s.c.NewYearConf.ActID)
	if err != nil {
		log.Error("s.dao.PrizeCardsCache(%d) err(%+v)", mid, err)
		err = nil
	}
	if len(res) != 0 {
		return
	}
	card1 := &model.PrizeCardRep{CardType: model.CardType1}
	card3 := &model.PrizeCardRep{CardType: model.CardType3}
	card12 := &model.PrizeCardRep{CardType: model.CardType12}
	userCards, err := s.dao.UserCards(c, mid, s.c.NewYearConf.ActID)
	if err != nil || userCards == nil {
		log.Error("s.dao.UserCards(%d) err(%+v)", mid, err)
		err = nil
		// 塞空数据
		res = append(res, card1, card3, card12)
		s.cache.Do(c, func(c context.Context) {
			if err = s.dao.SetPrizeCardsCache(context.Background(), mid, s.c.NewYearConf.ActID, res); err != nil {
				log.Error("s.dao.SetPrizeCardsCache(%d %+v) err(%+v)", mid, res, err)
			}
		})
		return
	}
	for _, uc := range userCards {
		switch uc.CardType {
		case model.CardType1:
			s.initCard(card1, uc)
		case model.CardType3:
			s.initCard(card3, uc)
		case model.CardType12:
			s.initCard(card12, uc)
		}
	}
	res = append(res, card1, card3, card12)
	s.cache.Do(c, func(c context.Context) {
		if err = s.dao.SetPrizeCardsCache(context.Background(), mid, s.c.NewYearConf.ActID, res); err != nil {
			log.Error("s.dao.SetPrizeCardsCache(%d %+v) err(%+v)", mid, res, err)
		}
	})
	return
}

func (s *Service) isDuring(c context.Context) (err error) {
	var now = time.Now().Unix()
	if s.c.NewYearConf.StartTime > now {
		err = ecode.CouponNewYearNotStartErr
		return
	}
	if s.c.NewYearConf.EndTime < now {
		err = ecode.CouponNewYearIsEndErr
		return
	}
	return
}

func (s *Service) initCard(card *model.PrizeCardRep, uc *model.CouponUserCard) {
	if uc.State >= model.CardStateOpened {
		card.State = uc.State
		FullAmount := model.MapFullAmount[uc.CardType]
		if bi, ok := s.allBranchInfo[uc.BatchToken]; ok {
			card.OriginalPrice = int64(FullAmount)
			card.CouponAmount = int64(bi.FullAmount - bi.Amount)
			card.DiscountRate = fmt.Sprintf("%.1f折", (bi.FullAmount-bi.Amount)*10/FullAmount)
		}
	}
}

// PrizeDraw .
func (s *Service) PrizeDraw(c context.Context, mid int64, cardType int8) (res *model.PrizeCardRep, err error) {
	if err = s.isDuring(c); err != nil {
		return
	}
	var (
		pcard *model.PrizeCardRep
		ucard *model.CouponUserCard
	)
	pcard, err = s.dao.PrizeCardCache(c, mid, s.c.NewYearConf.ActID, cardType)
	if err != nil {
		log.Error("s.dao.PrizeCardCache(%d %+v) err(%+v)", mid, pcard, err)
		err = nil
	}
	if pcard != nil && pcard.State != model.CardStateNotOpen {
		err = ecode.CouponNewYearIsOpenErr
		return
	}
	if pcard == nil {
		ucard, err = s.dao.UserCard(c, mid, s.c.NewYearConf.ActID, cardType)
		if err != nil {
			return
		}
		if ucard != nil {
			res = &model.PrizeCardRep{CardType: cardType}
			s.initCard(res, ucard)
			s.cache.Do(c, func(c context.Context) {
				if err = s.dao.SetPrizeCardCache(context.Background(), mid, s.c.NewYearConf.ActID, res); err != nil {
					log.Error("s.dao.SetPrizeCardCache(%d %+v) err(%+v)", mid, res, err)
				}
			})
			return
		}
	}
	vipInfo, err := s.vipinfoClient.Info(c, &v1.InfoReq{Mid: mid})
	if err != nil {
		log.Error("vipinfoSrv.Service.Info(%d) err(%+v)", mid, err)
		return
	}
	log.Info("vipInfo(%d %+v)", mid, vipInfo.Res)
	var batchToken = ""
	data := []byte(strconv.FormatInt(mid, 10) + model.CardSalt)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	md5Int, _ := strconv.ParseInt(md5str[24:], 16, 64)
	if md5Int%100 <= s.c.NewYearConf.RandNum {
		log.Info("good luck!(mid:%d)", mid)
		batchToken = s.MapNoVipBatchToken[cardType]
	} else {
		switch vipInfo.Res.Type {
		case vimdl.NotVip:
			// 历史从未充值过大会员
			if vipInfo.Res.DueDate == 0 {
				batchToken = s.MapNoVipBatchToken[cardType]
			} else {
				batchToken = s.MapMonthBatchToken[cardType]
			}
		case vimdl.Vip:
			batchToken = s.MapMonthBatchToken[cardType]
		case vimdl.AnnualVip:
			if (vipInfo.Res.DueDate/1000 - time.Now().Unix()) > 180*86400 {
				batchToken = s.MapMore180BatchToken[cardType]
			} else {
				batchToken = s.MapLess180BatchToken[cardType]
			}
		}
	}
	bi := s.allBranchInfo[batchToken]
	if bi == nil {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	if bi.State == model.BatchStateBlock {
		err = ecode.CouponBatchBlockErr
		return
	}
	pc, err := s.addPrizeCoupon(c, bi, mid, s.c.NewYearConf.ActID, cardType)
	if err != nil {
		log.Error("s.addPrizeCoupon(%d %+v, %d) err(%+v)", mid, bi, cardType, err)
		return
	}
	FullAmount := model.MapFullAmount[cardType]
	res = &model.PrizeCardRep{
		CardType:      cardType,
		State:         pc.State,
		OriginalPrice: int64(FullAmount),
		CouponAmount:  int64(bi.FullAmount - bi.Amount),
		DiscountRate:  fmt.Sprintf("%.1f折", (bi.FullAmount-bi.Amount)*10/FullAmount),
	}
	s.cache.Do(c, func(c context.Context) {
		if err = s.dao.SetPrizeCardCache(context.Background(), mid, s.c.NewYearConf.ActID, res); err != nil {
			log.Error("s.dao.SetPrizeCardCache(%d %+v) err(%+v)", mid, res, err)
		}
	})
	return
}

// UsableAllowanceCouponV2 usable coupon v2.
func (s *Service) UsableAllowanceCouponV2(c context.Context, a *coupinv1.UsableAllowanceCouponV2Req) (res *model.CouponTipInfo, err error) {
	var (
		us            []*model.CouponAllowancePanelInfo
		ui            []*model.CouponAllowancePanelInfo
		all           []*model.CouponAllowanceInfo
		inuse         []*model.CouponAllowanceInfo
		selectdPrice  float64
		currentPrice  float64
		maxAmount     float64
		bestCoupon    *model.CouponAllowancePanelInfo
		selectdCoupon *model.CouponAllowancePanelInfo
	)
	res = &model.CouponTipInfo{
		CouponTip: model.CouponTipNotUse,
	}
	if len(a.PriceInfo) == 0 {
		return
	}
	selectdPrice = a.PriceInfo[0].Price
	if all, err = s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   a.Mid,
		State: model.NotUsed,
	}); err != nil {
		return
	}
	// coupon using
	if inuse, err = s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
		Mid:   a.Mid,
		State: model.InUse,
	}); err != nil {
		return
	}
	all = append(all, inuse...)
	for _, v := range a.PriceInfo {
		if us, _, ui, err = s.UsableAllowanceCoupons(c,
			a.Mid,
			v.Price,
			all,
			int(v.Plat),
			int8(v.ProdLimMonth),
			int8(v.ProdLimRenewal)); err != nil {
			return
		}
		availables := append(us, ui...)
		if len(availables) == 0 {
			continue
		}
		sort.Slice(availables, func(i int, j int) bool {
			if availables[i].Amount == availables[j].Amount {
				return availables[i].Usable > availables[j].Usable
			}
			return availables[i].Amount > availables[j].Amount
		})
		if availables[0].Amount > maxAmount {
			currentPrice = v.Price
			maxAmount = availables[0].Amount
			bestCoupon = availables[0]
		}
		if v.Price == selectdPrice {
			selectdCoupon = availables[0]
			break
		}
	}
	if bestCoupon == nil {
		return
	}
	switch {
	case currentPrice == selectdPrice && bestCoupon.Usable == model.AllowanceUsable:
		res.CouponTip = fmt.Sprintf(model.CouponTipUse, bestCoupon.Amount)
		res.CouponInfo = bestCoupon
	case currentPrice == selectdPrice && bestCoupon.Usable == model.AllowanceDisables && bestCoupon.State == model.InUse:
		res.CouponTip = model.CouponTipInUse
	case selectdCoupon == nil:
		res.CouponTip = model.CouponTipChooseOther
	}
	return
}
