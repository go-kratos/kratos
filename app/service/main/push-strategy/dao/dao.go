package dao

import (
	"context"
	"time"

	"go-common/app/service/main/push-strategy/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	c                   *conf.Config
	db                  *xsql.DB
	redis               *redis.Pool
	mc                  *memcache.Pool
	httpClient          *xhttp.Client
	appsStmt            *xsql.Stmt
	businessesStmt      *xsql.Stmt
	addTaskStmt         *xsql.Stmt
	taskStmt            *xsql.Stmt
	settingsByRangeStmt *xsql.Stmt
	maxSettingIDStmt    *xsql.Stmt
	mcUUIDExpire        int32
	mcCDExpire          int32
	redisLimitDayExpire int32
}

// New new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                   c,
		db:                  xsql.NewMySQL(c.MySQL),
		httpClient:          xhttp.NewClient(c.HTTPClient),
		redis:               redis.NewPool(c.Redis.Config),
		mc:                  memcache.NewPool(c.Memcache.Config),
		mcUUIDExpire:        int32(time.Duration(c.Memcache.UUIDExpire) / time.Second),
		mcCDExpire:          int32(time.Duration(c.Memcache.CDExpire) / time.Second),
		redisLimitDayExpire: int32(time.Duration(c.Redis.LimitDayExpire) / time.Second),
	}
	d.appsStmt = d.db.Prepared(_appsSQL)
	d.businessesStmt = d.db.Prepared(_businessesSQL)
	d.addTaskStmt = d.db.Prepared(_addTaskSQL)
	d.taskStmt = d.db.Prepared(_taskSQL)
	d.settingsByRangeStmt = d.db.Prepared(_settingsByRangeSQL)
	d.maxSettingIDStmt = d.db.Prepared(_maxSettingIDSQL)
	return
}

// PromError prom error
func PromError(name string) {
	prom.BusinessErrCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	prom.BusinessInfoCount.Incr(name)
}

// BeginTx begin transaction.
func (d *Dao) BeginTx(ctx context.Context) (*xsql.Tx, error) {
	return d.db.Begin(ctx)
}

// Ping .
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	if err = d.pingMC(ctx); err != nil {
		return
	}
	return d.db.Ping(ctx)
}

// Close .
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}
