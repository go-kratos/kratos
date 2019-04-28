package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/up-rating/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao def dao struct
type Dao struct {
	c              *conf.Config
	db             *sql.DB
	rddb           *sql.DB
	redis          *redis.Pool
	upRatingExpire int64
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		db:             sql.NewMySQL(c.DB.Main),
		rddb:           sql.NewMySQL(c.DB.Slave),
		redis:          redis.NewPool(c.Redis.Config),
		upRatingExpire: int64(time.Duration(c.Redis.UpRatingExpire) / time.Second),
	}
	return d
}

// Ping ping db
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("d.db.Ping error(%v)", err)
		return
	}
	return
}

// Close close db conn
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}

// BeginTran begin transcation
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}
