package v1

import (
	"context"
	"time"

	"go-common/app/infra/config/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
)

// Dao dao.
type Dao struct {
	// mysql
	db *sql.DB
	// redis
	redis  *redis.Pool
	expire time.Duration
	// cache
	pathCache string
}

// New new a dao.
func New(c *conf.Config) *Dao {
	d := &Dao{
		// db
		db: sql.NewMySQL(c.DB),
		// redis
		redis:  redis.NewPool(c.Redis),
		expire: time.Duration(c.PollTimeout),
		// cache
		pathCache: c.PathCache,
	}
	return d
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}

// Close close resuouces.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}

// Ping ping is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	return d.db.Ping(c)
}
