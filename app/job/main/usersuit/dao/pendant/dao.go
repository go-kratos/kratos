package pendant

import (
	"context"

	"go-common/app/job/main/usersuit/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao struct info of Dao.
type Dao struct {
	db     *sql.DB
	c      *conf.Config
	client *bm.Client
	// redis
	redis *redis.Pool

	notifyURL string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     sql.NewMySQL(c.Mysql),
		client: bm.NewClient(c.HTTPClient),
		// redis
		redis:     redis.NewPool(c.PendantRedis.Config),
		notifyURL: c.NotifyURL,
	}

	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingRedis(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
