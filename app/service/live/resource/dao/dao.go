package dao

import (
	"context"
	"fmt"
	"go-common/library/net/http/blademaster"
	"math/rand"
	"time"

	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/lrucache"
	"go-common/library/database/orm"
	xsql "go-common/library/database/sql"

	"github.com/jinzhu/gorm"
)

const (
	_prefixResource  = "live:res:"
	_prefixLiveCheck = "live:check:"
	_prefixConf      = "live:conf:"
)

func cacheResourceKey(typ string, platform string, build int64) string {
	return fmt.Sprintf("%s:%s:%s:%d", _prefixResource, typ, platform, build)
}

func cacheLiveCheckKey(platform, system, mobile string) string {
	return fmt.Sprintf("%s:%s:%s:%s", _prefixLiveCheck, platform, system, mobile)
}

func cacheConfKey(key string) string {
	return fmt.Sprintf("%s:%s", _prefixConf, key)
}

// Dao dao
type Dao struct {
	c          *conf.Config
	db         *xsql.DB
	rsDB       *gorm.DB
	rsDBReader *gorm.DB
	sCache     []*lrucache.SyncCache
	client     *blademaster.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	if c.CacheCapacity <= 0 {
		c.CacheCapacity = 100
	}
	if c.CacheBucket <= 0 {
		c.CacheBucket = 100
	}
	if c.CacheInstCnt <= 0 {
		c.CacheInstCnt = 10
	}
	if c.CacheTimeout <= 0 {
		c.CacheTimeout = 60
	}
	sCache := make([]*lrucache.SyncCache, c.CacheInstCnt)
	for i := range sCache {
		sCache[i] = lrucache.NewSyncCache(c.CacheCapacity, c.CacheBucket, c.CacheTimeout)
	}
	dao = &Dao{
		c:          c,
		db:         xsql.NewMySQL(c.MySQL),
		rsDB:       orm.NewMySQL(c.DB.Resource),
		rsDBReader: orm.NewMySQL(c.DB.ResourceReader),
		sCache:     sCache,
		client:     blademaster.NewClient(c.HttpClient),
	}
	rand.Seed(time.Now().UnixNano())
	return
}

// mysql log
func (d *Dao) initORM() {
	d.rsDB.LogMode(true)
	d.rsDB.LogMode(true)
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
	d.rsDB.Close()
	d.rsDBReader.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	// TODO: if you need use mc,redis, please add
	if err = d.db.Ping(c); err != nil {
		return
	}
	return
	//return d.pingRedis(c)
}
