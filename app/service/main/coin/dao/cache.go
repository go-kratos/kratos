package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func countKey(mid int64) string {
	return "uct2_" + strconv.FormatInt(mid, 10)
}

func (d *Dao) cacheSFUserCoin(mid int64) string {
	return "uct_" + strconv.FormatInt(mid, 10)
}

func itemCoinKey(aid, tp int64) string {
	return fmt.Sprintf("itc_%d_%d", tp, aid)
}

func (d *Dao) cacheSFItemCoin(aid int64, tp int64) string {
	return fmt.Sprintf("itc_%d_%d", aid, tp)
}

func expKey(mid int64) (key string) {
	now := time.Now()
	key = fmt.Sprintf("exp_%d_%02d%02d", mid, now.Month(), now.Day())
	return
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// get user coin count.
	// mc: -key=countKey
	CacheUserCoin(c context.Context, mid int64) (count float64, err error)
	// set user coin count
	// mc: -key=countKey -expire=d.mcExpire
	AddCacheUserCoin(c context.Context, mid int64, count float64) (err error)
	// mc: -key=itemCoinKey
	CacheItemCoin(c context.Context, aid int64, tp int64) (count int64, err error)
	// mc: -key=itemCoinKey -expire=d.mcExpire
	AddCacheItemCoin(c context.Context, aid int64, count int64, tp int64) (err error)
	// mc: -key=expKey -type=get
	Exp(c context.Context, mid int64) (exp int64, err error)
	// mc: -key=expKey -expire=d.expireExp
	SetTodayExpCache(c context.Context, mid int64, exp int64) (err error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=-1 -check_null_code=$==-1 -singleflight=true
	UserCoin(c context.Context, mid int64) (count float64, err error)
	// cache: -nullcache=-1 -check_null_code=$==-1 -singleflight=true
	ItemCoin(c context.Context, aid int64, tp int64) (count int64, err error)
}
