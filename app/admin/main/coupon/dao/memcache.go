package dao

import (
	"context"
	"fmt"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixCouponAllowances = "cas:%d:%d"
	_prefixCoupons          = "cs:%d:%d"
	_prefixGrantUnique      = "gu:%s"
)

func couponAllowancesKey(mid int64, state int8) string {
	return fmt.Sprintf(_prefixCouponAllowances, mid, state)
}

func couponsKey(mid int64, ct int8) string {
	return fmt.Sprintf(_prefixCoupons, ct, mid)
}

func grantUnique(token string) string {
	return fmt.Sprintf(_prefixGrantUnique, token)
}

// DelCouponAllowancesKey delete allowances cache.
func (d *Dao) DelCouponAllowancesKey(c context.Context, mid int64, state int8) (err error) {
	return d.delCache(c, couponAllowancesKey(mid, state))
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
		}
	}
	return
}

//DelCouponTypeCache del coupon.
func (d *Dao) DelCouponTypeCache(c context.Context, mid int64, ct int8) (err error) {
	return d.delCache(c, couponsKey(mid, ct))
}

//DelGrantUniqueLock del lock.
func (d *Dao) DelGrantUniqueLock(c context.Context, token string) (err error) {
	return d.delCache(c, grantUnique(token))
}

// AddGrantUniqueLock add grant coupon use lock.
func (d *Dao) AddGrantUniqueLock(c context.Context, token string, seconds int32) (succeed bool) {
	var (
		key  = grantUnique(token)
		conn = d.mc.Get(c)
		err  error
	)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: seconds,
	}
	if err = conn.Add(item); err != nil {
		if err != gmc.ErrNotStored {
			log.Error("mc.Add(%s) error(%v)", key, err)
		}
		return
	}
	succeed = true
	return
}
