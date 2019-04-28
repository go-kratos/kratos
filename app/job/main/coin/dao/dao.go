package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/coin/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/stat/prom"
)

// Dao define dao.
type Dao struct {
	c                   *conf.Config
	coinDB              *sql.DB
	hitSettlePeriodStmt *sql.Stmt
	getSettlePeriodStmt *sql.Stmt
	getTotalCoinsStmt   []*sql.Stmt
	redis               *redis.Pool
	loginExpire         int32
}

// PromError .
func PromError(name string) {
	prom.BusinessErrCount.Incr(name)
}

// New new and return service.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                 c,
		coinDB:            sql.NewMySQL(c.DB.Coin),
		getTotalCoinsStmt: make([]*sql.Stmt, SHARDING),
		redis:             redis.NewPool(c.Redis),
		loginExpire:       int32(time.Duration(c.CoinJob.LoginExpire) / time.Second),
	}
	for i := 0; i < SHARDING; i++ {
		d.getTotalCoinsStmt[i] = d.coinDB.Prepared(fmt.Sprintf(_getTotalCoins, i))
	}
	d.hitSettlePeriodStmt = d.coinDB.Prepared(_hitSettlePeriod)
	d.getSettlePeriodStmt = d.coinDB.Prepared(_getSettlePeriod)
	return
}

// Ping check service health.
func (d *Dao) Ping(c context.Context) error {
	return d.coinDB.Ping(c)
}

// Close close sevice.
func (d *Dao) Close() {
	d.coinDB.Close()
}
