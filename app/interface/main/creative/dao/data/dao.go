package data

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/data"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"go-common/library/database/hbase.v2"
)

// Dao is data dao.
type Dao struct {
	c *conf.Config
	// http
	client     *bm.Client
	statClient *bm.Client
	// uri
	statURI     string
	tagV2URI    string
	coverBFSURI string
	// mc
	mc          *memcache.Pool
	mcExpire    int32
	mcIdxExpire int32
	statCacheOn bool
	// hbase
	hbase        *hbase.Client
	hbaseTimeOut time.Duration
	// chan
	missch chan func()
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// http
		client:      bm.NewClient(c.HTTPClient.Slow),
		statClient:  bm.NewClient(c.HTTPClient.Slow),
		statURI:     c.Host.Data + _statURL,
		tagV2URI:    c.Host.Tag + _tagv2URL,
		coverBFSURI: c.Host.Coverrec + _coverURL,
		// memcache
		mc:          memcache.NewPool(c.Memcache.Data.Config),
		mcExpire:    int32(time.Duration(c.Memcache.Data.DataExpire) / time.Second),
		mcIdxExpire: int32(time.Duration(c.Memcache.Data.IndexExpire) / time.Second),
		statCacheOn: c.StatCacheOn,
		// hbase
		hbase: hbase.NewClient(c.HBase.Config),
		// chan
		missch:       make(chan func(), 1024),
		hbaseTimeOut: time.Duration(time.Millisecond * 200),
	}
	go d.cacheproc()
	return
}

// Stat get user stat play/fans/...
func (d *Dao) Stat(c context.Context, ip string, mid int64) (st *data.Stat, err error) {
	// try cache
	if st, _ = d.statCache(c, mid); st != nil {
		return
	}
	// from api
	if st, err = d.stat(c, mid, ip); st != nil {
		d.AddCache(func() {
			d.addStatCache(context.TODO(), mid, st)
		})
	}
	return
}

// UpStat get up stat from hbase
func (d *Dao) UpStat(c context.Context, mid int64, dt string) (st *data.UpBaseStat, err error) {
	// try cache
	if d.statCacheOn {
		if st, _ = d.upBaseStatCache(c, mid, dt); st != nil {
			log.Info("upBaseStatCache cache found mid(%d)", mid)
			return
		}
	}
	// from hbase
	if st, err = d.BaseUpStat(c, mid, dt); st != nil {
		d.AddCache(func() {
			d.addUpBaseStatCache(context.TODO(), mid, dt, st)
		})
	}
	return
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("mc.ping.Store error(%v)", err)
	}
	return
}

// Close mc close
func (d *Dao) Close() (err error) {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.hbase != nil {
		d.hbase.Close()
	}
	return
}

// AddCache add to chan for cache
func (d *Dao) AddCache(f func()) {
	select {
	case d.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (d *Dao) cacheproc() {
	for {
		f := <-d.missch
		f()
	}
}
