package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/dm/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

var (
	missedCount = prom.CacheMiss
	cachedCount = prom.CacheHit
)

// Dao dao struct
type Dao struct {
	conf       *conf.Config
	httpClient *bm.Client
	// redis
	redisDM          *redis.Pool
	redisDMIDExpire  int32
	redisLockExpire  int32
	redisIndexExpire int32
	redisVideoExpire int32
	// database
	dmMetaReader *sql.DB // bilibili_dm_meta
	biliDM       *sql.DB // blibili_dm
	dmWriter     *sql.DB
	// databus
	databus *databus.Databus
	// elastic
	es *elastic.Elastic
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:       c,
		httpClient: bm.NewClient(c.HTTPClient),
		// redis
		redisDM:          redis.NewPool(c.Redis.DM.Config),
		redisDMIDExpire:  int32(time.Duration(c.Redis.DM.DMIDExpire) / time.Second),
		redisLockExpire:  int32(time.Duration(c.Redis.DM.LockExpire) / time.Second),
		redisIndexExpire: int32(time.Duration(c.Redis.DM.IndexExpire) / time.Second),
		redisVideoExpire: int32(time.Duration(c.Redis.DM.VideoExpire) / time.Second),
		// database
		dmMetaReader: sql.NewMySQL(c.DB.DMMetaReader),
		biliDM:       sql.NewMySQL(c.DB.DM),
		dmWriter:     sql.NewMySQL(c.DB.DMWriter),
		// databus
		databus: databus.New(c.Databus),
		// elastic
		es: elastic.NewElastic(c.ES),
	}
	return
}

// PromCacheHit prom cache hit
func PromCacheHit(name string) {
	cachedCount.Incr(name)
}

// PromCacheMiss prom cache hit
func PromCacheMiss(name string) {
	missedCount.Incr(name)
}

// Ping ping dao status
func (d *Dao) Ping(c context.Context) (err error) {
	// dm redis
	rds := d.redisDM.Get(c)
	defer rds.Close()
	if _, err = rds.Do("SET", "ping", "pong"); err != nil {
		log.Error("rds.Do(SET) error(%v)", err)
		return
	}
	// database
	if err = d.dmMetaReader.Ping(c); err != nil {
		log.Error("dmMetaReader.Ping() error(%v)", err)
		return
	}
	if err = d.biliDM.Ping(c); err != nil {
		log.Error("biliDM.Ping() error(%v)", err)
		return
	}
	return
}
