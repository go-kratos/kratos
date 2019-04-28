package dao

import (
	"context"
	"net/http"
	"time"

	"go-common/app/service/main/passport-sns/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
)

// Dao dao struct
type Dao struct {
	c        *conf.Config
	db       *sql.DB
	mc       *memcache.Pool
	client   *http.Client
	mcExpire int32
}

// New create new dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		db:       sql.NewMySQL(c.MySQL),
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		client: &http.Client{
			Timeout: 600 * time.Millisecond,
		},
	}
	return
}

// Ping check db and mc health.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}
