package dao

import (
	"context"
	"time"

	"go-common/app/service/main/dynamic/conf"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

const (
	_regionURL    = "/dynamic/region"
	_regionTagURL = "/dynamic/tag"
	_liveURL      = "/room/v1/Area/dynamic"
	_hotURL       = "/x/internal/tag/hotmap"
	_pridURL      = "/x/internal/tag/prids"
)

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Dao dao.
type Dao struct {
	// http
	httpR *bm.Client
	// bigData api
	regionURI    string
	regionTagURI string
	// live api
	liveURI string
	// tag api
	hotURI  string
	pridURI string
	// memcache
	mc       *memcache.Pool
	mcExpire int32
	// cache Prom
	cacheProm *prom.Prom
}

// New dao new.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		httpR:        bm.NewClient(c.HTTPClient.Read),
		regionURI:    c.Host.BigDataURI + _regionURL,
		regionTagURI: c.Host.BigDataURI + _regionTagURL,
		liveURI:      c.Host.LiveURI + _liveURL,
		hotURI:       c.Host.APIURI + _hotURL,
		pridURI:      c.Host.APIURI + _pridURL,
		// memcache
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	d.cacheProm = prom.CacheHit
	return
}

// Ping check connection success.
func (dao *Dao) Ping(c context.Context) (err error) {
	err = dao.pingMC(c)
	return
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.mc != nil {
		dao.mc.Close()
	}
}
