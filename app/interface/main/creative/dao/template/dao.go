package template

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

// Dao is archive dao.
type Dao struct {
	// config
	c *conf.Config
	// db
	db *sql.DB
	// insert
	addTplStmt *sql.Stmt
	// update
	upTplStmt  *sql.Stmt
	delTplStmt *sql.Stmt
	// select
	getTplStmt      *sql.Stmt
	getMutilTplStmt *sql.Stmt
	getCntStmt      *sql.Stmt
	// mc
	mc       *memcache.Pool
	mcExpire int32
	// chan
	tplch chan func()
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Creative),
		// memcache
		mc:       memcache.NewPool(c.Memcache.Archive.Config),
		mcExpire: int32(time.Duration(c.Memcache.Archive.TplExpire) / time.Second),
		// chan
		tplch: make(chan func(), 1024),
	}
	// insert
	d.addTplStmt = d.db.Prepared(_addTplSQL)
	// update
	d.upTplStmt = d.db.Prepared(_upTplSQL)
	d.delTplStmt = d.db.Prepared(_delTplSQL)
	// select
	d.getTplStmt = d.db.Prepared(_getTplSQL)
	d.getMutilTplStmt = d.db.Prepared(_getMutilTplSQL)
	d.getCntStmt = d.db.Prepared(_getCntSQL)
	go d.cacheproc()
	return
}

// Ping db
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	return d.db.Ping(c)
}

// Close db
func (d *Dao) Close() (err error) {
	if d.db != nil {
		d.db.Close()
	}
	return d.db.Close()
}

// addCache add to chan for cache
func (d *Dao) addCache(f func()) {
	select {
	case d.tplch <- f:
	default:
		log.Warn("template cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (d *Dao) cacheproc() {
	for {
		f := <-d.tplch
		f()
	}
}
