package dao

import (
	"context"
	"math/rand"
	"time"

	"go-common/app/service/main/spy/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao event dao def.
type Dao struct {
	c *conf.Config
	// db
	db *sql.DB
	// mc
	mcUser       *memcache.Pool
	mcUserExpire int32
	// redis
	redis        *redis.Pool
	expire       int
	verifyExpire int
	// http client
	httpClient   *bm.Client
	auditInfoURI string
	// databus for spy score change
	dbScoreChange *databus.Databus
	r             *rand.Rand
}

// New create instance of dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// db
		db: sql.NewMySQL(c.DB.Spy),
		// mc
		mcUser:       memcache.NewPool(c.Memcache.User),
		mcUserExpire: int32(time.Duration(c.Memcache.UserExpire) / time.Second),
		// redis
		redis:        redis.NewPool(c.Redis.Config),
		expire:       int(time.Duration(c.Redis.Expire) / time.Second),
		verifyExpire: int(time.Duration(c.Redis.VerifyCdTimes) / time.Second),
		// http client
		httpClient:   bm.NewClient(c.HTTPClient),
		auditInfoURI: c.Account + _auditInfo,
		// databus
		dbScoreChange: databus.New(c.DBScoreChange),
		r:             rand.New(rand.NewSource(time.Now().Unix())),
	}
	if conf.Conf.Property.UserInfoShard <= 0 {
		panic("conf.Conf.Property.UserInfoShard <= 0")
	}
	if conf.Conf.Property.HistoryShard <= 0 {
		panic("conf.Conf.Property.HistoryShard <= 0")
	}
	return
}

// Ping check db health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		return
	}
	if err = d.PingRedis(c); err != nil {
		return
	}
	return d.db.Ping(c)
}

// Close close all db connections.
func (d *Dao) Close() {
	if d.mcUser != nil {
		d.mcUser.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}
