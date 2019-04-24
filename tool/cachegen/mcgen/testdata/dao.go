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
	articleExpire int32
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
		articleExpire: int32(5),
	}
	return
}

//go:generate kratos tool mcgen
type _mc interface {
	// mc: -key=articleKey
	CacheArticles(c context.Context, keys []int64) (map[int64]*Article, error)
	// mc: -key=articleKey
	CacheArticle(c context.Context, key int64) (*Article, error)
	// mc: -key=keyMid
	CacheArticle1(c context.Context, key int64, mid int64) (*Article, error)
	// mc: -key=noneKey
	CacheNone(c context.Context) (*Article, error)
	// mc: -key=articleKey
	CacheString(c context.Context, key int64) (string, error)

	// mc: -key=articleKey -expire=d.articleExpire -encode=json
	AddCacheArticles(c context.Context, values map[int64]*Article) error
	// mc: -key=article2Key -expire=d.articleExpire -encode=json
	AddCacheArticles2(c context.Context, values map[int64]*Article, tp int64) error
	// 这里也支持自定义注释 会替换默认的注释
	// mc: -key=articleKey -expire=d.articleExpire -encode=json|gzip
	AddCacheArticle(c context.Context, key int64, value *Article) error
	// mc: -key=keyMid -expire=d.articleExpire -encode=gob
	AddCacheArticle1(c context.Context, key int64, value *Article, mid int64) error
	// mc: -key=noneKey
	AddCacheNone(c context.Context, value *Article) error
	// mc: -key=articleKey -expire=d.articleExpire
	AddCacheString(c context.Context, key int64, value string) error

	// mc: -key=articleKey
	DelCacheArticles(c context.Context, keys []int64) error
	// mc: -key=articleKey
	DelCacheArticle(c context.Context, key int64) error
	// mc: -key=keyMid
	DelCacheArticle1(c context.Context, key int64, mid int64) error
	// mc: -key=noneKey
	DelCacheNone(c context.Context) error
}

func articleKey(id int64) string {
	return fmt.Sprintf("art_%d", id)
}

func article2Key(id, tp int64) string {
	return fmt.Sprintf("art_%d_%d", id, tp)
}

func keyMid(id, mid int64) string {
	return fmt.Sprintf("art_%d_%d", id, mid)
}

func noneKey() string {
	return "none"
}
