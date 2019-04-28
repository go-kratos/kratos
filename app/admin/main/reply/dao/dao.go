package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	es "go-common/library/database/elastic"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao dao.
type Dao struct {
	c *conf.Config
	// db
	db      *sql.DB
	dbSlave *sql.DB
	// cache
	redis       *redis.Pool
	redisExpire int32
	mc          *memcache.Pool
	mcExpire    int32
	// http
	httpClient *bm.Client
	// databus
	eventBus *databus.Databus
	// new databus stats
	statsBus   *databus.Databus
	statsTypes map[int32]string
	es         *es.Elastic
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// db
		db:      sql.NewMySQL(c.DB.Reply),
		dbSlave: sql.NewMySQL(c.DB.ReplySlave),
		// cache
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
		mc:          memcache.NewPool(c.Memcache.Config),
		mcExpire:    int32(time.Duration(c.Memcache.Expire) / time.Second),
		// http
		httpClient: bm.NewClient(c.HTTPClient),
		// databus
		eventBus: databus.New(c.Databus.Event),
		es:       es.NewElastic(c.Es),
	}
	// new databus stats
	d.statsTypes = make(map[int32]string)
	for name, typ := range c.StatTypes {
		d.statsTypes[typ] = name
	}
	d.statsBus = databus.New(c.Databus.Stats)
	return
}

func hit(id int64) int64 {
	return id % 200
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// Ping ping resouces is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	return d.pingMC(c)
}
