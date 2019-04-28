package block

import (
	"context"
	"math"
	"math/rand"
	"time"

	"go-common/app/service/main/member/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

type notifyFunc func(context.Context, int64, string) error

// Dao is
type Dao struct {
	*cacheTTL
	c                *conf.Config
	mc               *memcache.Pool
	db               *sql.DB
	client           *bm.Client
	NotifyPurgeCache notifyFunc
}

type cacheTTL struct {
	UserTTL     int32
	UserMaxRate float64
	UserT       float64
}

// New is
func New(conf *conf.Config, db *sql.DB, mc *memcache.Pool, client *bm.Client, notifyFunc notifyFunc) *Dao {
	d := &Dao{
		c:                conf,
		mc:               mc,
		db:               db,
		client:           client,
		NotifyPurgeCache: notifyFunc,
	}
	d.cacheTTL = newCacheTTL(conf.BlockCacheTTL)
	return d
}

// BeginTran is
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	if tx, err = d.db.Begin(c); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func durationToSeconds(expire xtime.Duration) int32 {
	return int32(time.Duration(expire) / time.Second)
}

func newCacheTTL(c *conf.BlockCacheTTL) *cacheTTL {
	return &cacheTTL{
		UserTTL:     durationToSeconds(c.UserTTL),
		UserMaxRate: c.UserMaxRate,
		UserT:       c.UserT,
	}
}

func (ttl *cacheTTL) mcUserExpire(key string) (sec int32) {
	if ttl.UserT == 0.0 {
		return ttl.UserTTL
	}
	// rate = -log(1-x)/t
	rate := -math.Log(1-rand.Float64()) / ttl.UserT
	if rate <= 1.0 {
		return ttl.UserTTL
	}
	if rate > ttl.UserMaxRate {
		rate = ttl.UserMaxRate
	}
	sec = int32(rate * float64(ttl.UserTTL))
	if rate >= 5.0 {
		log.Info("mc hotkey : %s, expire rate : %.2f , time : %d", key, rate, sec)
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}
