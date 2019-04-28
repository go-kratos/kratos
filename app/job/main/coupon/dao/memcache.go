package dao

import (
	"context"
	"fmt"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixCoupons      = "cs:%d:%d"
	_prefixCouponBlance = "cbl:%d:%d"
	_prefixprizeCard    = "nypc:%d:%d:%d" // 元旦活动卡片
	_prefixprizeCards   = "nypcs:%d:%d"   // 元旦活动卡片列表
)

func couponBalancesKey(mid int64, ct int8) string {
	return fmt.Sprintf(_prefixCouponBlance, ct, mid)
}

// DelCouponBalancesCache delete user coupons blance cache.
func (d *Dao) DelCouponBalancesCache(c context.Context, mid int64, ct int8) (err error) {
	return d.delCache(c, couponBalancesKey(mid, ct))
}

func couponsKey(mid int64, ct int8) string {
	return fmt.Sprintf(_prefixCoupons, ct, mid)
}

func prizeCardKey(mid, actID int64, cardType int8) string {
	return fmt.Sprintf(_prefixprizeCard, mid, actID, cardType)
}

func prizeCardsKey(mid, actID int64) string {
	return fmt.Sprintf(_prefixprizeCards, mid, actID)
}

// DelCouponsCache delete user coupons cache.
func (d *Dao) DelCouponsCache(c context.Context, mid int64, ct int8) (err error) {
	return d.delCache(c, couponsKey(mid, ct))
}

// DelCache del cache.
func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			log.Warn("delCache ErrNotFound(%s)", key)
			err = nil
		} else {
			err = errors.Wrapf(err, "mc.Delete(%s)", key)
		}
	}
	return
}

// DelPrizeCardKey .
func (d *Dao) DelPrizeCardKey(c context.Context, mid, actID int64, cardType int8) (err error) {
	return d.delCache(c, prizeCardKey(mid, actID, cardType))
}

// DelPrizeCardsKey .
func (d *Dao) DelPrizeCardsKey(c context.Context, mid, actID int64) (err error) {
	log.Warn("DelPrizeCardsKey(%s)", prizeCardsKey(mid, actID))
	return d.delCache(c, prizeCardsKey(mid, actID))
}
