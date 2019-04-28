package dao

import (
	"context"

	"go-common/app/service/bbq/recsys-recall/conf"
	xcache "go-common/app/service/bbq/recsys-recall/dao/cache"
	"go-common/library/cache/redis"
)

// Dao dao
type Dao struct {
	c       *conf.Config
	redis   *redis.Pool
	bfredis *redis.Pool
	lcache  *xcache.LocalCache
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:       c,
		redis:   redis.NewPool(c.Redis),
		bfredis: redis.NewPool(c.BFRedis),
		lcache:  xcache.NewLocalCache(c.LocalCache),
	}
	return
}

// GetInvertedIndex 获取倒排索引
func (d *Dao) GetInvertedIndex(ctx context.Context, key string, force bool) (b []byte, err error) {
	if b = d.lcache.Get(key); b != nil && !force {
		return
	}

	conn := d.redis.Get(ctx)
	defer conn.Close()

	for retry := 0; retry < 3; retry++ {
		b, err = redis.Bytes(conn.Do("GET", key))
		if err == redis.ErrNil {
			b = make([]byte, 0)
			return
		}
		if b != nil {
			d.lcache.Set(key, b)
			return
		}
	}

	return
}

// SetInvertedIndex 更新倒排索引
func (d *Dao) SetInvertedIndex(c context.Context, key string, value []byte) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	_, err := conn.Do("SETEX", key, 86400, value)

	return err
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return nil
}
