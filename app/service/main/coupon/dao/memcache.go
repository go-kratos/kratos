package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/coupon/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixCoupons          = "cs:%d:%d"
	_prefixUseUnique        = "cu:%s:%d"
	_prefixCouponBlance     = "cbl:%d:%d"
	_prefixGrantUnique      = "gr:%s:%d"
	_prefixBranchCount      = "bcu:%s"
	_useLockTimeout         = 10
	_prefixCouponAllowances = "cas:%d:%d"
	_receiveLog             = "rl:%s%s%d"
	_uniqueNo               = "uq:%s"
	_useUniqueNoTimeout     = 1296000         //15天
	_prefixprizeCard        = "nypc:%d:%d:%d" // 元旦活动卡片
	_prefixprizeCards       = "nypcs:%d:%d"   // 元旦活动卡片列表
)

func receiveLogKey(appkey, orderNo string, ct int8) string {
	return fmt.Sprintf(_receiveLog, appkey, orderNo, ct)
}

func couponsKey(mid int64, ct int8) string {
	return fmt.Sprintf(_prefixCoupons, ct, mid)
}

func useUniqueKey(orderNO string, ct int8) string {
	return fmt.Sprintf(_prefixUseUnique, orderNO, ct)
}

func couponBalancesKey(mid int64, ct int8) string {
	return fmt.Sprintf(_prefixCouponBlance, ct, mid)
}

func userGrantKey(token string, mid int64) string {
	return fmt.Sprintf(_prefixGrantUnique, token, mid)
}

func branchCurrentCount(token string) string {
	return fmt.Sprintf(_prefixBranchCount, token)
}

func couponAllowancesKey(mid int64, state int8) string {
	return fmt.Sprintf(_prefixCouponAllowances, mid, state)
}

func prizeCardKey(mid, actID int64, cardType int8) string {
	return fmt.Sprintf(_prefixprizeCard, mid, actID, cardType)
}

func prizeCardsKey(mid, actID int64) string {
	return fmt.Sprintf(_prefixprizeCards, mid, actID)
}
func couponuniqueNoKey(uniqueno string) string {
	return fmt.Sprintf(_uniqueNo, uniqueno)
}

// DelUniqueKey delete  use coupon lock cache.
func (d *Dao) DelUniqueKey(c context.Context, orderNO string, ct int8) (err error) {
	return d.delCache(c, useUniqueKey(orderNO, ct))
}

// DelCouponsCache delete user coupons cache.
func (d *Dao) DelCouponsCache(c context.Context, mid int64, ct int8) (err error) {
	return d.delCache(c, couponsKey(mid, ct))
}

// DelCouponBalancesCache delete user coupons blance cache.
func (d *Dao) DelCouponBalancesCache(c context.Context, mid int64, ct int8) (err error) {
	return d.delCache(c, couponBalancesKey(mid, ct))
}

// DelGrantKey delete  user grant lock cache.
func (d *Dao) DelGrantKey(c context.Context, token string, mid int64) (err error) {
	return d.delCache(c, userGrantKey(token, mid))
}

// DelBranchCurrentCountKey delete branch current cache.
func (d *Dao) DelBranchCurrentCountKey(c context.Context, token string) (err error) {
	return d.delCache(c, branchCurrentCount(token))
}

// DelCouponAllowancesKey delete allowances cache.
func (d *Dao) DelCouponAllowancesKey(c context.Context, mid int64, state int8) (err error) {
	return d.delCache(c, couponAllowancesKey(mid, state))
}

// DelPrizeCardKey .
func (d *Dao) DelPrizeCardKey(c context.Context, mid, actID int64, cardType int8) (err error) {
	return d.delCache(c, prizeCardKey(mid, actID, cardType))
}

// DelPrizeCardsKey .
func (d *Dao) DelPrizeCardsKey(c context.Context, mid, actID int64) (err error) {
	return d.delCache(c, prizeCardsKey(mid, actID))
}

// CouponsCache coupons cache.
func (d *Dao) CouponsCache(c context.Context, mid int64, ct int8) (coupons []*model.CouponInfo, err error) {
	var (
		key  = couponsKey(mid, ct)
		item *gmc.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err = conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "mc.Get(%s)", key)
		d.errProm.Incr("get_mc")
		return
	}
	couponInfoList := &model.PointInfoList{}
	if err = conn.Scan(item, couponInfoList); err != nil {
		err = errors.Wrapf(err, "mc.Scan(%s)", key)
		d.errProm.Incr("scan_mc")
		return
	}
	coupons = couponInfoList.PointInfoList
	if coupons == nil {
		coupons = []*model.CouponInfo{}
	}
	return
}

// SetCouponsCache set coupons cache.
func (d *Dao) SetCouponsCache(c context.Context, mid int64, ct int8, coupons []*model.CouponInfo) (err error) {
	var (
		expire = d.mcExpire
		key    = couponsKey(mid, ct)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: &model.PointInfoList{PointInfoList: coupons}, Expiration: expire, Flags: gmc.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "mc.Set(%s)", key)
		d.errProm.Incr("set_mc")
	}
	return
}

// AddUseUniqueLock add coupon use lock.
func (d *Dao) AddUseUniqueLock(c context.Context, orderNO string, ct int8) (succeed bool) {
	var (
		key  = useUniqueKey(orderNO, ct)
		conn = d.mc.Get(c)
		err  error
	)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: _useLockTimeout,
	}
	if err = conn.Add(item); err != nil {
		if err != gmc.ErrNotStored {
			log.Error("mc.Add(%s) error(%v)", key, err)
			d.errProm.Incr("add_mc")
		}
	} else {
		succeed = true
	}
	return
}

// AddReceiveUniqueLock add coupon use lock.
func (d *Dao) AddReceiveUniqueLock(c context.Context, appkey, orderNO string, ct int8) (succeed bool) {
	var (
		key  = receiveLogKey(appkey, orderNO, ct)
		conn = d.mc.Get(c)
		err  error
	)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: _useLockTimeout,
	}
	if err = conn.Add(item); err != nil {
		if err != gmc.ErrNotStored {
			log.Error("mc.Add(%s) error(%v)", key, err)
			d.errProm.Incr("add_mc")
		}
	} else {
		succeed = true
	}
	return
}

//DelReceiveUniqueLock del receive lock.
func (d *Dao) DelReceiveUniqueLock(c context.Context, appkey, orderNO string, ct int8) (err error) {
	err = d.delCache(c, receiveLogKey(appkey, orderNO, ct))
	return
}

// DelCache del cache.
func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "mc.Delete(%s)", key)
			d.errProm.Incr("del_mc")
		}
	}
	return
}

// CouponBlanceCache coupon blance cache.
func (d *Dao) CouponBlanceCache(c context.Context, mid int64, ct int8) (coupons []*model.CouponBalanceInfo, err error) {
	var (
		key  = couponBalancesKey(mid, ct)
		item *gmc.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err = conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "mc.Get(%s)", key)
		d.errProm.Incr("get_mc")
		return
	}
	couponBlanceList := &model.CouponBalanceList{}
	if err = conn.Scan(item, couponBlanceList); err != nil {
		err = errors.Wrapf(err, "mc.Scan(%s)", key)
		d.errProm.Incr("scan_mc")
		return
	}
	coupons = couponBlanceList.CouponBalanceList
	if coupons == nil {
		coupons = []*model.CouponBalanceInfo{}
	}
	return
}

// SetCouponBlanceCache set coupon blance cache.
func (d *Dao) SetCouponBlanceCache(c context.Context, mid int64, ct int8, coupons []*model.CouponBalanceInfo) (err error) {
	var (
		expire = d.mcExpire
		key    = couponBalancesKey(mid, ct)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: &model.CouponBalanceList{CouponBalanceList: coupons}, Expiration: expire, Flags: gmc.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "mc.Set(%s)", key)
		d.errProm.Incr("set_mc")
	}
	return
}

// AddUniqueNoLock add grant coupon use lock.
func (d *Dao) AddUniqueNoLock(c context.Context, uniqueno string) (succeed bool) {
	var (
		key  = couponuniqueNoKey(uniqueno)
		conn = d.mc.Get(c)
		err  error
	)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: _useUniqueNoTimeout,
	}
	if err = conn.Add(item); err != nil {
		if err != gmc.ErrNotStored {
			log.Error("mc.Add(%s) error(%v)", key, err)
			d.errProm.Incr("add_mc")
		}
	} else {
		succeed = true
	}
	return
}

// AddGrantUniqueLock add grant unique coupon use lock.
func (d *Dao) AddGrantUniqueLock(c context.Context, token string, mid int64) (succeed bool) {
	var (
		key  = userGrantKey(token, mid)
		conn = d.mc.Get(c)
		err  error
	)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: _useLockTimeout,
	}
	if err = conn.Add(item); err != nil {
		if err != gmc.ErrNotStored {
			log.Error("mc.Add(%s) error(%v)", key, err)
			d.errProm.Incr("add_mc")
		}
	} else {
		succeed = true
	}
	return
}

//BranchCurrentCountCache branchInfo current count cache.
func (d *Dao) BranchCurrentCountCache(c context.Context, token string) (count int, err error) {
	var (
		key  = branchCurrentCount(token)
		conn = d.mc.Get(c)
		item *gmc.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			count = -1
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	if err = conn.Scan(item, &count); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%+v)", item)
		return
	}
	return
}

// SetBranchCurrentCountCache set branch current cache.
func (d *Dao) SetBranchCurrentCountCache(c context.Context, token string, count int) (err error) {
	var (
		key  = branchCurrentCount(token)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&gmc.Item{Key: key, Object: count, Flags: gmc.FlagJSON, Expiration: d.mcExpire}); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%+v)", key, count)
		return
	}
	return
}

// IncreaseBranchCurrentCountCache increase branch current count cache.
func (d *Dao) IncreaseBranchCurrentCountCache(c context.Context, token string, count uint64) (err error) {
	var (
		key  = branchCurrentCount(token)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Increment(key, count); err != nil {
		err = errors.Wrapf(err, "conn.Increment(%s,%d)", key, count)
		return
	}
	return
}

// CouponAllowanceCache coupon allowance cache.
func (d *Dao) CouponAllowanceCache(c context.Context, mid int64, state int8) (coupons []*model.CouponAllowanceInfo, err error) {
	var (
		key  = couponAllowancesKey(mid, state)
		item *gmc.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err = conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			coupons = nil
			return
		}
		err = errors.Wrapf(err, "mc.Get(%s)", key)
		d.errProm.Incr("get_mc")
		return
	}
	couponAllowanceList := &model.CouponAllowanceList{}
	if err = conn.Scan(item, couponAllowanceList); err != nil {
		err = errors.Wrapf(err, "mc.Scan(%s)", key)
		d.errProm.Incr("scan_mc")
		return
	}
	coupons = couponAllowanceList.CouponAllowanceList
	if coupons == nil {
		coupons = []*model.CouponAllowanceInfo{}
	}
	return
}

// SetCouponAllowanceCache set coupon allowance cache.
func (d *Dao) SetCouponAllowanceCache(c context.Context, mid int64, state int8, coupons []*model.CouponAllowanceInfo) (err error) {
	var (
		expire = d.mcExpire
		key    = couponAllowancesKey(mid, state)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: &model.CouponAllowanceList{CouponAllowanceList: coupons}, Expiration: expire, Flags: gmc.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "mc.Set(%s)", key)
		d.errProm.Incr("set_mc")
	}
	return
}

// SetPrizeCardCache .
func (d *Dao) SetPrizeCardCache(c context.Context, mid, actID int64, prizeCard *model.PrizeCardRep) (err error) {
	var (
		expire = d.prizeExpire
		key    = prizeCardKey(mid, actID, prizeCard.CardType)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: prizeCard, Expiration: expire, Flags: gmc.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "mc.Set(%s)", key)
		d.errProm.Incr("set_mc")
	}
	return
}

// SetPrizeCardsCache .
func (d *Dao) SetPrizeCardsCache(c context.Context, mid, actID int64, prizeCards []*model.PrizeCardRep) (err error) {
	var (
		expire = d.prizeExpire
		key    = prizeCardsKey(mid, actID)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: &model.PrizeCards{List: prizeCards}, Expiration: expire, Flags: gmc.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "mc.Set(%s)", key)
		d.errProm.Incr("set_mc")
	}
	return
}

// PrizeCardCache .
func (d *Dao) PrizeCardCache(c context.Context, mid, actID int64, cardType int8) (prizeCard *model.PrizeCardRep, err error) {
	var (
		key  = prizeCardKey(mid, actID, cardType)
		item *gmc.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err = conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "mc.Get(%s)", key)
		d.errProm.Incr("get_mc")
		return
	}
	prizeCard = &model.PrizeCardRep{}
	if err = conn.Scan(item, &prizeCard); err != nil {
		err = errors.Wrapf(err, "mc.Scan(%s)", key)
		d.errProm.Incr("scan_mc")
		return
	}
	return
}

// PrizeCardsCache .
func (d *Dao) PrizeCardsCache(c context.Context, mid, actID int64) (prizeCards []*model.PrizeCardRep, err error) {
	var (
		key  = prizeCardsKey(mid, actID)
		item *gmc.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err = conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "mc.Get(%s)", key)
		d.errProm.Incr("get_mc")
		return
	}
	PrizeCardlist := &model.PrizeCards{}
	if err = conn.Scan(item, PrizeCardlist); err != nil {
		err = errors.Wrapf(err, "mc.Scan(%s)", key)
		d.errProm.Incr("scan_mc")
		return
	}
	prizeCards = PrizeCardlist.List
	if prizeCards == nil {
		prizeCards = []*model.PrizeCardRep{}
	}
	return
}
