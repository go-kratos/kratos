package share

import (
	"context"

	"go-common/app/interface/main/web-goblin/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Dao dao struct.
type Dao struct {
	// config
	c *conf.Config
	// db
	db *sql.DB
	// redis
	redis *redis.Pool
}

// New new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c:     c,
		db:    sql.NewMySQL(c.DB.Goblin),
		redis: redis.NewPool(c.Redis.Config),
	}
	return
}

// Ping ping dao
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	return
}

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}
