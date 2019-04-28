package dao

import (
	"context"
	"time"

	"go-common/app/job/main/tag/conf"
	filgrpc "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao dao struct .
type Dao struct {
	maxNum    int
	expNewArc int
	conf      *conf.Config
	redis     *redis.Pool
	redisRank *redis.Pool
	redisTag  *redis.Pool
	memcache  *memcache.Pool
	platform  *sql.DB
	// grpc
	filClient filgrpc.FilterClient
}

// New init Dao .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:      c,
		platform:  sql.NewMySQL(c.Platform.MySQL),
		maxNum:    c.Tag.MaxArcsLimit,
		redis:     redis.NewPool(c.Redis.Rank.Redis),
		redisRank: redis.NewPool(c.RedisRank),
		redisTag:  redis.NewPool(c.RedisTag),
		memcache:  memcache.NewPool(c.Memcache),
		expNewArc: int(time.Duration(c.Redis.Rank.Expire.TagNewArc) / time.Second),
	}
	var err error
	if d.filClient, err = filgrpc.NewClient(c.FilterGRPC); err != nil {
		panic(err)
	}
	return
}

// PingRedis ping redis .
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// pingRedisRank ping rank redis .
func (d *Dao) pingRedisRank(c context.Context) (err error) {
	conn := d.redisRank.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// pingRedisTag ping rank redis .
func (d *Dao) pingRedisTag(c context.Context) (err error) {
	conn := d.redisTag.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// Ping ping db .
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.platform.Ping(c); err != nil {
		log.Error("ping db error(%v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		log.Error("ping redis error(%v)", err)
		return
	}
	if err = d.pingRedisRank(c); err != nil {
		log.Error("ping rank redis error(%v)", err)
		return
	}
	if err = d.pingRedisTag(c); err != nil {
		log.Error("ping tag redis error(%v)", err)
	}
	return
}

// Close close .
func (d *Dao) Close() (err error) {
	if d.redis != nil {
		d.redis.Close()
	}
	if d.platform != nil {
		d.platform.Close()
	}
	if d.redisRank != nil {
		d.redisRank.Close()
	}
	return
}

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.platform.Begin(c)
}
