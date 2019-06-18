package testdata

import (
	"context"
	"fmt"
	"time"

	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/container/pool"
	xtime "github.com/bilibili/kratos/pkg/time"
)

// Dao .
type Dao struct {
	mc            *memcache.Memcache
	demoExpire int32
}

// New new dao
func New() (d *Dao) {
	cfg := &memcache.Config{
		Config: &pool.Config{
			Active:      10,
			Idle:        5,
			IdleTimeout: xtime.Duration(time.Second),
		},
		Name:  "test",
		Proto: "tcp",
		// Addr:         "172.16.33.54:11214",
		Addr:         "127.0.0.1:11211",
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	d = &Dao{
		mc:            memcache.New(cfg),
		demoExpire: int32(5),
	}
	return
}

//go:generate kratos tool genmc
type _mc interface {
	// mc: -key=demoKey
	CacheDemos(c context.Context, keys []int64) (map[int64]*Demo, error)
	// mc: -key=demoKey
	CacheDemo(c context.Context, key int64) (*Demo, error)
	// mc: -key=keyMid
	CacheDemo1(c context.Context, key int64, mid int64) (*Demo, error)
	// mc: -key=noneKey
	CacheNone(c context.Context) (*Demo, error)
	// mc: -key=demoKey
	CacheString(c context.Context, key int64) (string, error)

	// mc: -key=demoKey -expire=d.demoExpire -encode=json
	AddCacheDemos(c context.Context, values map[int64]*Demo) error
	// mc: -key=demo2Key -expire=d.demoExpire -encode=json
	AddCacheDemos2(c context.Context, values map[int64]*Demo, tp int64) error
	// 这里也支持自定义注释 会替换默认的注释
	// mc: -key=demoKey -expire=d.demoExpire -encode=json|gzip
	AddCacheDemo(c context.Context, key int64, value *Demo) error
	// mc: -key=keyMid -expire=d.demoExpire -encode=gob
	AddCacheDemo1(c context.Context, key int64, value *Demo, mid int64) error
	// mc: -key=noneKey
	AddCacheNone(c context.Context, value *Demo) error
	// mc: -key=demoKey -expire=d.demoExpire
	AddCacheString(c context.Context, key int64, value string) error

	// mc: -key=demoKey
	DelCacheDemos(c context.Context, keys []int64) error
	// mc: -key=demoKey
	DelCacheDemo(c context.Context, key int64) error
	// mc: -key=keyMid
	DelCacheDemo1(c context.Context, key int64, mid int64) error
	// mc: -key=noneKey
	DelCacheNone(c context.Context) error
}

func demoKey(id int64) string {
	return fmt.Sprintf("art_%d", id)
}

func demo2Key(id, tp int64) string {
	return fmt.Sprintf("art_%d_%d", id, tp)
}

func keyMid(id, mid int64) string {
	return fmt.Sprintf("art_%d_%d", id, mid)
}

func noneKey() string {
	return "none"
}
