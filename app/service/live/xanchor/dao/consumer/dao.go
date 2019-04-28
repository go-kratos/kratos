package consumer

import (
	"context"

	"go-common/app/service/live/xanchor/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

type refValue struct {
	field string
	v     interface{}
}

// Dao dao
type Dao struct {
	c         *conf.Config
	redis     *redis.Pool
	db        *xsql.DB
	dbLiveApp *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:         c,
		redis:     redis.NewPool(c.Redis),
		db:        xsql.NewMySQL(c.MySQL),
		dbLiveApp: xsql.NewMySQL(c.LiveAppMySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}
