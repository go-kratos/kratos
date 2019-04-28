package dao

import (
	"context"

	"go-common/app/service/main/share/conf"
	"go-common/library/cache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Dao is redis dao.
type Dao struct {
	c   *conf.Config
	db  *sql.DB
	rds *redis.Pool
	// cache chan
	cacheCh chan func()
	// databus
	databus *databus.Databus
	// archiveDatabus
	archiveDatabus *databus.Databus
	// sources
	sources map[int64]struct{}
	// asyncCache
	asyncCache *cache.Cache
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		db:             sql.NewMySQL(c.DB),
		rds:            redis.NewPool(c.Redis),
		cacheCh:        make(chan func(), 1024),
		databus:        databus.New(c.Databus),
		archiveDatabus: databus.New(c.ArchiveDatabus),
		sources:        make(map[int64]struct{}, len(c.Sources)),
		asyncCache:     cache.New(c.Worker, 1024),
	}
	for _, s := range c.Sources {
		d.sources[s] = struct{}{}
	}
	return d
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}

// Close close
func (d *Dao) Close() {
	d.db.Close()
	d.rds.Close()
}

// Ping ping
func (d *Dao) Ping() (err error) {
	if err = d.db.Ping(context.Background()); err != nil {
		log.Error("d.db.Ping error(%v)", err)
		return
	}
	conn := d.rds.Get(context.Background())
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("redis.Set error(%v)", err)
		return
	}
	return
}
