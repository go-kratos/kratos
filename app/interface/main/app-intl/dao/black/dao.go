package black

import (
	"context"
	"time"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/library/cache/redis"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao is black dao.
type Dao struct {
	// http clientAsyn
	clientAsyn *httpx.Client
	// redis
	redis     *redis.Pool
	expireRds int32
	aCh       chan func()
}

// New new a black dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		clientAsyn: httpx.NewClient(c.HTTPClientAsyn),
		// redis init
		redis:     redis.NewPool(c.Redis.Feed.Config),
		expireRds: int32(time.Duration(c.Redis.Feed.ExpireBlack) / time.Second),
		aCh:       make(chan func(), 1024),
	}
	go d.cacheproc()
	return
}

// Ping  Ping check redis connection
func (d *Dao) Ping(c context.Context) (err error) {
	connRedis := d.redis.Get(c)
	_, err = connRedis.Do("SET", "PING", "PONG")
	connRedis.Close()
	return
}

// AddBlacklist is.
func (d *Dao) AddBlacklist(mid, aid int64) {
	d.addCache(func() {
		d.addBlackCache(context.Background(), mid, aid)
	})
}

// DelBlacklist is.
func (d *Dao) DelBlacklist(mid, aid int64) {
	d.addCache(func() {
		d.delBlackCache(context.Background(), mid, aid)
	})
}

// BlackList is.
func (d *Dao) BlackList(c context.Context, mid int64) (aidm map[int64]struct{}, err error) {
	var ok bool
	if ok, err = d.expireBlackCache(c, mid); err != nil {
		return
	}
	if ok {
		aidm, err = d.blackCache(c, mid)
	}
	return
}

// addCache add cache to mc by goroutine
func (d *Dao) addCache(i func()) {
	select {
	case d.aCh <- i:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc cache proc
func (d *Dao) cacheproc() {
	for {
		f, ok := <-d.aCh
		if !ok {
			log.Warn("cache proc exit")
			return
		}
		f()
	}
}
