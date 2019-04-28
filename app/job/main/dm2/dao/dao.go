package dao

import (
	"context"
	"time"

	"go-common/app/job/main/dm2/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/bfs"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_pageSize = 1000
)

// Dao dao struct
type Dao struct {
	conf *conf.Config
	// batch query size
	pageSize int
	// database
	dmWriter     *sql.DB
	dmReader     *sql.DB
	biliDMWriter *sql.DB
	// redis
	dmRds       *redis.Pool
	dmRdsExpire int32
	// recent dm redis
	dmRctRds    *redis.Pool
	dmRctExpire int32
	// segment dm redis
	dmSegRds    *redis.Pool
	dmSegExpire int32
	// memcache
	mc               *memcache.Pool
	mcExpire         int32
	subtitleMc       *memcache.Pool
	subtitleMcExpire int32
	// memcache new
	dmSegMC       *memcache.Pool
	dmSegMCExpire int32
	// recent dm redis
	rctRds       *redis.Pool
	rctRdsExpire int32

	// http client
	httpCli *bm.Client
	// upload dm
	bfsCli *bfs.BFS
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:             c,
		dmWriter:         sql.NewMySQL(c.DB.DMWriter),
		dmReader:         sql.NewMySQL(c.DB.DMReader),
		biliDMWriter:     sql.NewMySQL(c.DB.BiliDMWriter),
		dmRds:            redis.NewPool(c.Redis.DM.Config),
		dmRdsExpire:      int32(time.Duration(c.Redis.DM.Expire) / time.Second),
		dmRctRds:         redis.NewPool(c.Redis.DMRct.Config),
		dmRctExpire:      int32(time.Duration(c.Redis.DMRct.Expire) / time.Second),
		dmSegRds:         redis.NewPool(c.Redis.DMSeg.Config),
		dmSegExpire:      int32(time.Duration(c.Redis.DMSeg.Expire) / time.Second),
		mc:               memcache.NewPool(c.Memcache.Config),
		mcExpire:         int32(time.Duration(c.Memcache.Expire) / time.Second),
		subtitleMc:       memcache.NewPool(c.SubtitleMemcache.Config),
		subtitleMcExpire: int32(time.Duration(c.SubtitleMemcache.Expire) / time.Second),
		dmSegMC:          memcache.NewPool(c.DMMemcache.Config),
		dmSegMCExpire:    int32(time.Duration(c.DMMemcache.Expire) / time.Second),
		rctRds:           redis.NewPool(c.Redis.DMRct.Config),
		rctRdsExpire:     int32(time.Duration(c.Redis.DMRct.Expire) / time.Second),
		httpCli:          bm.NewClient(c.HTTPClient),
		bfsCli:           bfs.New(c.Bfs.Client),
		pageSize:         int(c.DB.QueryPageSize),
	}
	if d.pageSize <= 0 {
		d.pageSize = _pageSize
	}
	return
}

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.dmWriter.Begin(c)
}

// BeginBiliDMTran .
func (d *Dao) BeginBiliDMTran(c context.Context) (*sql.Tx, error) {
	return d.biliDMWriter.Begin(c)
}

// Ping dm dao ping.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.dmWriter.Ping(c); err != nil {
		log.Error("dmWriter.Ping() error(%v)", err)
		return
	}
	if err = d.dmReader.Ping(c); err != nil {
		log.Error("dmReader.Ping() error(%v)", err)
		return
	}
	if err = d.biliDMWriter.Ping(c); err != nil {
		log.Error("biliDMWriter.Ping() error(%v)", err)
		return
	}
	// mc
	mconn := d.mc.Get(c)
	defer mconn.Close()
	if err = mconn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("mc.Set error(%v)", err)
		return
	}
	// dm redis
	dmRdsConn := d.dmRds.Get(c)
	defer dmRdsConn.Close()
	if _, err = dmRdsConn.Do("SET", "ping", "pong"); err != nil {
		log.Error("dmRds.Set error(%v)", err)
		return
	}
	rctRdsConn := d.dmRctRds.Get(c)
	defer rctRdsConn.Close()
	if _, err = rctRdsConn.Do("SET", "ping", "pong"); err != nil {
		log.Error("rctRds.Set error(%v)", err)
		return
	}
	dmSegConn := d.dmSegRds.Get(c)
	defer dmSegConn.Close()
	if _, err = dmSegConn.Do("SET", "ping", "pong"); err != nil {
		log.Error("dmSegConn.Set error(%v)", err)
		return
	}
	return
}
