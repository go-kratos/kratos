package dao

import (
	"context"
	"time"

	"go-common/app/job/main/dm/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_subjectSharding = 100
	_indexSharding   = 1000
)

// Dao dao struct.
type Dao struct {
	// redis
	redis       *redis.Pool
	redisExpire int32
	// memcache
	mc       *memcache.Pool
	mcExpire int32
	// mysql
	dmReader *sql.DB
	dmWriter *sql.DB
}

// New return dm dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
		// memcache
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		// mysql
		dmReader: sql.NewMySQL(c.DB.DMReader),
		dmWriter: sql.NewMySQL(c.DB.DMWriter),
	}
	return
}

func (d *Dao) hitSubject(oid int64) int64 {
	return oid % _subjectSharding
}

func (d *Dao) hitIndex(oid int64) int64 {
	return oid % _indexSharding
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
	// mc
	mconn := d.mc.Get(c)
	defer mconn.Close()
	if err = mconn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("mc.Set error(%v)", err)
		return
	}
	// dm redis
	rconn := d.redis.Get(c)
	defer rconn.Close()
	if _, err = rconn.Do("SET", "ping", "pong"); err != nil {
		rconn.Close()
		log.Error("redis.Set error(%v)", err)
	}
	return
}
