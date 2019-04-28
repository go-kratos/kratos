package dao

import (
	"context"
	"time"

	"go-common/app/interface/live/push-live/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao dao
type Dao struct {
	c                         *conf.Config
	db                        *xsql.DB
	httpClient                *xhttp.Client
	relationHBase             *hbase.Client
	relationHBaseReadTimeout  time.Duration
	blackListHBase            *hbase.Client
	blackListHBaseReadTimeout time.Duration
}

// Prometheus
var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                         c,
		db:                        xsql.NewMySQL(c.MySQL),
		relationHBase:             hbase.NewClient(c.HBase.Config),
		relationHBaseReadTimeout:  time.Duration(c.HBase.ReadTimeout),
		httpClient:                xhttp.NewClient(c.HTTPClient),
		blackListHBase:            hbase.NewClient(c.BlackListHBase.Config),
		blackListHBaseReadTimeout: time.Duration(c.BlackListHBase.ReadTimeout),
	}
	return
}

// RedisOption return redis options
func (d *Dao) RedisOption() []redis.DialOption {
	cnop := redis.DialConnectTimeout(time.Duration(d.c.Redis.PushInterval.DialTimeout))
	rdop := redis.DialReadTimeout(time.Duration(d.c.Redis.PushInterval.ReadTimeout))
	wrop := redis.DialWriteTimeout(time.Duration(d.c.Redis.PushInterval.WriteTimeout))
	return []redis.DialOption{cnop, rdop, wrop}
}

// Close close the resource.
func (d *Dao) Close() (err error) {
	if err = d.relationHBase.Close(); err != nil {
		log.Error("[dao.dao|Close] d.relationHBase.Close() error(%v)", err)
		PromError("hbase:close")
	}
	if err = d.db.Close(); err != nil {
		log.Error("[dao.dao|Close] d.db.Close() error(%v)", err)
		PromError("db:close")
	}
	return
}

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}

//PromInfoAdd add prom info by value
func PromInfoAdd(name string, value int64) {
	infosCount.Add(name, value)
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		PromError("mysql:Ping")
		log.Error("[dao.dao|Ping] d.db.Ping error(%v)", err)
		return
	}
	if err = d.relationHBase.Ping(c); err != nil {
		PromError("hbase:Ping")
		log.Error("[dao.dao|Ping] d.relationHBase.Ping error(%v)", err)
		return
	}
	return
}
