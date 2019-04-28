package history

import (
	"runtime"
	"time"

	hisrpc "go-common/app/interface/main/history/rpc/client"
	"go-common/app/interface/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Dao is account dao.
type Dao struct {
	// rpc
	hisRPC    *hisrpc.Service
	mc        *memcache.Pool
	mCh       chan func()
	expireHis int32
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		hisRPC:    hisrpc.New(c.HisRPC),
		mc:        memcache.NewPool(c.Memcache.Config),
		mCh:       make(chan func(), 10240),
		expireHis: int32(time.Duration(c.Memcache.HisExpire) / time.Second),
	}
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go d.cacheproc()
	}
	return
}

var (
	cachedCount = prom.CacheHit
	missedCount = prom.CacheMiss
)

// addCache add archive to mc or redis
func (d *Dao) addCache(f func()) {
	select {
	case d.mCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc write memcache and stat redis use goroutine
func (d *Dao) cacheproc() {
	for {
		f := <-d.mCh
		f()
	}
}
