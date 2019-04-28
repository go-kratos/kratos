package dao

import (
	"context"
	"time"

	"go-common/app/service/main/coin/conf"
	"go-common/app/service/main/coin/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// Dao dao config.
type Dao struct {
	c *conf.Config
	// databus stat
	stat *databus.Databus
	// coin
	coin *sql.DB
	//http
	httpClient *bm.Client
	// databus
	dbBigData *databus.Databus
	// databus for coin-job
	dbCoinJob *databus.Databus
	// redis
	redis       *redis.Pool
	expireAdded int32
	expireExp   int32
	// tag url
	tagURI        string
	mc            *memcache.Pool
	mcExpire      int32
	cache         *fanout.Fanout
	Businesses    map[int64]*model.Business
	BusinessNames map[string]*model.Business
	es            *elastic.Elastic
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		coin:          sql.NewMySQL(c.DB.Coin),
		httpClient:    bm.NewClient(c.HTTPClient),
		tagURI:        c.TagURL,
		redis:         redis.NewPool(c.Redis.Config),
		mc:            memcache.NewPool(c.Memcache.Config),
		mcExpire:      int32(time.Duration(c.Memcache.Expire) / time.Second),
		expireExp:     int32(time.Duration(c.Memcache.ExpExpire) / time.Second),
		dbBigData:     databus.New(c.DbBigData),
		dbCoinJob:     databus.New(c.DbCoinJob),
		stat:          databus.New(c.Stat.Databus),
		expireAdded:   int32(time.Duration(c.Redis.Expire) / time.Second),
		cache:         fanout.New("cache", fanout.Buffer(10240)),
		Businesses:    make(map[int64]*model.Business),
		BusinessNames: make(map[string]*model.Business),
		es:            elastic.NewElastic(nil),
	}
	if len(c.Businesses) > 0 {
		for _, b := range c.Businesses {
			d.Businesses[b.ID] = b
			d.BusinessNames[b.Name] = b
		}
	}
	return
}

// PromError .
func PromError(name string) {
	prom.BusinessErrCount.Incr(name)
}

// Ping check service health.
func (dao *Dao) Ping(c context.Context) (err error) {
	return dao.coin.Ping(c)
}
