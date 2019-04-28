package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/tag/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_shard = 200
)

// Dao struct .
type Dao struct {
	conf               *conf.Config
	db                 *sql.DB
	redis              *redis.Pool
	redisExpire        int32
	subExpire          int32
	actionExpire       int32
	operateExpire      int32
	mc                 *memcache.Pool
	tagExpire          int32
	resExpire          int32
	resAllTidExpire    int32
	channelGroupExpire int32
}

// New dao return struct .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:               c,
		db:                 sql.NewMySQL(c.MySQL),
		redis:              redis.NewPool(c.Redis.Config),
		redisExpire:        int32(time.Duration(c.Redis.Expire) / time.Second),
		subExpire:          int32(time.Duration(c.Redis.SubExpire) / time.Second),
		actionExpire:       int32(time.Duration(c.Redis.ActionExpire) / time.Second),
		operateExpire:      int32(time.Duration(c.Redis.OperateExpire) / time.Second),
		mc:                 memcache.NewPool(c.Memcache.Config),
		tagExpire:          int32(time.Duration(c.Memcache.TagExpire) / time.Second),
		resExpire:          int32(time.Duration(c.Memcache.ResExpire) / time.Second),
		resAllTidExpire:    int32(time.Duration(c.Memcache.ResAllTidsExpire) / time.Second),
		channelGroupExpire: int32(time.Duration(c.Memcache.ChannelGroupExpire) / time.Second),
	}
	return
}

func (d *Dao) hit(mid int64) string {
	return fmt.Sprintf("%03d", mid%int64(_shard))
}

func (d *Dao) batchKey(oids []int64) (res map[string][]int64) {
	res = make(map[string][]int64, len(oids))
	for _, oid := range oids {
		key := fmt.Sprintf("%03d", oid%int64(_shard))
		res[key] = append(res[key], oid)
	}
	return
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	if tx, err = d.db.Begin(c); err != nil {
		log.Error("d.db.Begin error(%v)", err)
	}
	return
}

// Ping ping connection is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	if err = d.pingRedis(c); err != nil {
		return
	}
	err = d.pingMC(c)
	return
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.tagExpire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// Close close .
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}
