package dao

import (
	"context"
	"time"

	"go-common/app/service/main/favorite/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao favorite dao.
type Dao struct {
	db          *sql.DB
	dbRead      *sql.DB
	dbPush      *sql.DB
	mc          *memcache.Pool
	redis       *redis.Pool
	redisExpire int
	mcExpire    int32
	jobDatabus  *databus.Databus
	httpClient  *httpx.Client
}

// New a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// db
		db:     sql.NewMySQL(c.MySQL.Fav),
		dbRead: sql.NewMySQL(c.MySQL.Read),
		dbPush: sql.NewMySQL(c.MySQL.Push),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int(time.Duration(c.Redis.Expire) / time.Second),
		// memcache
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		// databus
		jobDatabus: databus.New(c.JobDatabus),
		// httpclient
		httpClient: httpx.NewClient(c.HTTPClient),
	}
	return
}

// Close close all connection.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.dbRead != nil {
		d.dbRead.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	if d.jobDatabus != nil {
		d.jobDatabus.Close()
	}
}

// BeginTran crate a *sql.Tx for database transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// Ping check connection used in dao
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	if err = d.pingMC(c); err != nil {
		return
	}
	err = d.pingMySQL(c)
	return
}
