package dao

import (
	"time"

	"go-common/app/service/main/secure/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"go-common/library/database/hbase.v2"
)

// Dao struct info of Dao.
type Dao struct {
	db                *sql.DB
	ddldb             *sql.DB
	c                 *conf.Config
	redis             *redis.Pool
	hbase             *hbase.Client
	locsExpire        int32
	expire            int64
	doubleCheckExpire int64
	mc                *memcache.Pool
	// http
	httpClient *bm.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                 c,
		redis:             redis.NewPool(c.Redis.Config),
		expire:            int64(time.Duration(c.Redis.Expire) / time.Second),
		doubleCheckExpire: int64(time.Duration(c.Redis.DoubleCheck) / time.Second),
		db:                sql.NewMySQL(c.Mysql.Secure),
		ddldb:             sql.NewMySQL(c.Mysql.DDL),
		hbase:             hbase.NewClient(c.HBase.Config),
		mc:                memcache.NewPool(c.Memcache.Config),
		locsExpire:        int32(time.Duration(c.Memcache.Expire) / time.Second),
		httpClient:        bm.NewClient(c.HTTPClient),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
