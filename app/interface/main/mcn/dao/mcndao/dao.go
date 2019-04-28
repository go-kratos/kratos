package mcndao

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/dao/global"
	"go-common/library/cache/memcache"
	"go-common/library/sync/pipeline/fanout"

	"github.com/bluele/gcache"
	"github.com/jinzhu/gorm"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	mc    *memcache.Pool
	mcndb *gorm.DB
	//cache tool
	cache         *fanout.Fanout
	mcnSignExpire int32
	mcnDataExpire int32
	localcache    gcache.Cache
}

// New init mysql db
func New(c *conf.Config, localcache gcache.Cache) (dao *Dao) {
	dao = &Dao{
		c:  c,
		mc: global.GetMc(),
		// cache worker
		cache:         fanout.New("cache", fanout.Worker(runtime.NumCPU()), fanout.Buffer(1024)),
		mcnSignExpire: int32(time.Duration(c.Memcache.McnSignCacheExpire) / time.Second),
		mcnDataExpire: int32(time.Duration(c.Memcache.McnDataCacheExpire) / time.Second),
		localcache:    localcache,
	}
	if localcache == nil {
		dao.localcache = gcache.New(c.RankCache.Size).Simple().Build()
	}
	if dao.mcnDataExpire == 0 {
		dao.mcnDataExpire = 3600
	}
	if dao.mcnSignExpire == 0 {
		dao.mcnSignExpire = 3600
	}
	var err error
	dao.mcndb, err = gorm.Open("mysql", c.MCNorm.DSN)
	if err != nil {
		panic(fmt.Errorf("db connect fail, err=%s", err))
	}
	dao.mcndb.LogMode(c.Other.Debug)
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.cache.Close()
	d.mc.Close()
	d.mcndb.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}

//GetMcnDB get mcn db
func (d *Dao) GetMcnDB() *gorm.DB {
	return d.mcndb
}
