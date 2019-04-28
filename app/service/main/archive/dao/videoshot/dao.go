package videoshot

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"go-common/app/service/main/archive/conf"
	"go-common/app/service/main/archive/model/videoshot"
)

// Dao is videoshot dao.
type Dao struct {
	// mysql
	db      *sql.DB
	dbRead  *sql.DB
	getStmt *sql.Stmt
	inStmt  *sql.Stmt
	// redis
	rds *redis.Pool
	// prom
	infoProm *prom.Prom
	// chan
	cacheCh chan func()
}

// New new a videoshot dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:      sql.NewMySQL(c.DB.Arc),
		dbRead:  sql.NewMySQL(c.DB.ArcRead),
		rds:     redis.NewPool(c.Redis.Archive.Config),
		cacheCh: make(chan func(), 1024),
	}
	d.getStmt = d.dbRead.Prepared(_getSQL)
	d.inStmt = d.db.Prepared(_inSQL)
	d.infoProm = prom.BusinessInfoCount
	go d.cacheproc()
	return d
}

// Videoshot get videoshot.
func (d *Dao) Videoshot(c context.Context, cid int64) (v *videoshot.Videoshot, err error) {
	var count, ver int
	if count, ver, err = d.cache(c, cid); err != nil {
		log.Error("d.cache(%d) error(%v)", cid, err)
		err = nil // NOTE: ignore error use db
	}
	if count != 0 {
		v = &videoshot.Videoshot{Cid: cid, Count: count}
		v.SetVersion(ver)
		return
	}
	if v, err = d.videoshot(c, cid); err != nil || v == nil {
		log.Warn("d.videoshot(%d) error(%v) or v==nil", cid, err)
		return
	}
	d.cacheCh <- func() {
		d.addCache(context.TODO(), v.Cid, v.Version(), v.Count)
	}
	return
}

// AddVideoshot add videoshot.
func (d *Dao) AddVideoshot(c context.Context, v *videoshot.Videoshot) (err error) {
	if _, err = d.addVideoshot(c, v); err != nil {
		log.Error("d.addVideoshot(%v) error(%v)", v, err)
		return
	}
	d.cacheCh <- func() {
		d.addCache(context.TODO(), v.Cid, v.Version(), v.Count)
	}
	return
}

// Close close resource.
func (d *Dao) Close() {
	if d.rds != nil {
		d.rds.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
	close(d.cacheCh)
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.rds.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

func (d *Dao) cacheproc() {
	for {
		f, ok := <-d.cacheCh
		if !ok {
			return
		}
		f()
	}
}
