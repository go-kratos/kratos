package dao

import (
	"context"
	"time"

	"go-common/app/job/main/thumbup/conf"
	"go-common/library/cache/memcache"
	xredis "go-common/library/cache/redis"
	"go-common/library/database/tidb"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Dao .
type Dao struct {
	// config
	c *conf.Config
	// tidb
	tidb     *tidb.DB
	itemTidb *tidb.DB
	// memcache
	mc            *memcache.Pool
	mcStatsExpire int32
	// redis
	redis                *xredis.Pool
	redisStatsExpire     int64
	redisUserLikesExpire int64
	redisItemLikesExpire int64
	// databus
	statDbus *databus.Databus
	// stmt
	businessesStmt      *tidb.Stmts
	itemLikesStmt       *tidb.Stmts
	userLikesStmt       *tidb.Stmts
	likeStateStmt       *tidb.Stmts
	statStmt            *tidb.Stmts
	updateLikeStateStmt *tidb.Stmts
}

// New .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// tidb
		tidb:     tidb.NewTiDB(c.Tidb),
		itemTidb: tidb.NewTiDB(c.ItemTidb),
		// memcache
		mc:            memcache.NewPool(c.Memcache.Config),
		mcStatsExpire: int32(time.Duration(c.Memcache.StatsExpire) / time.Second),
		// redis
		redis:                xredis.NewPool(c.Redis.Config),
		redisStatsExpire:     int64(time.Duration(c.Redis.StatsExpire) / time.Second),
		redisUserLikesExpire: int64(time.Duration(c.Redis.UserLikesExpire) / time.Second),
		redisItemLikesExpire: int64(time.Duration(c.Redis.ItemLikesExpire) / time.Second),
		// databus
		statDbus: databus.New(c.Databus.Stat),
	}

	d.statStmt = d.tidb.Prepared(_statSQL)
	d.itemLikesStmt = d.itemTidb.Prepared(_itemLikesSQL)
	d.businessesStmt = d.tidb.Prepared(_businessesSQL)
	d.userLikesStmt = d.tidb.Prepared(_userLikeListSQL)
	d.likeStateStmt = d.tidb.Prepared(_likeStateSQL)
	d.updateLikeStateStmt = d.tidb.Prepared(_updateLikeStateSQL)

	return
}

// Close .
func (d *Dao) Close() {
	d.mc.Close()
	d.tidb.Close()
	d.redis.Close()
}

// Ping .
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		log.Error("d.pingMC error(%v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		log.Error("d.pingRedis error(%v)", err)
	}
	return
}

// pingMC ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 100}
	err = conn.Set(&item)
	return
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("SET", "PING", "PONG")
	return
}
