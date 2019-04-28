package dao

import (
	"context"

	"go-common/app/admin/main/laser/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"time"
)

// Dao struct.
type Dao struct {
	c          *conf.Config
	laserDB    *sql.DB
	mc         *memcache.Pool
	mcExpire   int32
	HTTPClient *bm.Client
}

// New Dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		laserDB:    sql.NewMySQL(c.Mysql),
		mc:         memcache.NewPool(c.Memcache.Laser.Config),
		mcExpire:   int32(time.Duration(c.Memcache.Laser.Expire) / time.Second),
		HTTPClient: bm.NewClient(c.HTTPClient),
	}
	return
}

// Ping check db connection.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.laserDB.Ping(c)
}

// Close dao resources.
func (d *Dao) Close(c context.Context) (err error) {
	return d.laserDB.Close()
}
