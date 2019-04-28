package dao

import (
	"context"
	"time"

	"go-common/app/admin/ep/saga/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao def
type Dao struct {
	// cache
	httpClient     *bm.Client
	db             *gorm.DB
	mc             *memcache.Pool
	redis          *redis.Pool
	mcRecordExpire int32
}

// New create instance of Dao
func New() (d *Dao) {
	d = &Dao{
		mc:             memcache.NewPool(conf.Conf.Memcache.MC),
		httpClient:     bm.NewClient(conf.Conf.HTTPClient),
		db:             orm.NewMySQL(conf.Conf.ORM),
		redis:          redis.NewPool(conf.Conf.Redis),
		mcRecordExpire: int32(time.Duration(conf.Conf.Memcache.MCRecordExpire) / time.Second),
	}
	return
}

// Ping dao.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	if err = d.pingMC(c); err != nil {
		return
	}
	return d.db.DB().Ping()
}

// Close dao.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}
