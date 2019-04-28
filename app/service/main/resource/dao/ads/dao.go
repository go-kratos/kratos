package ads

import (
	"context"
	"time"

	"go-common/app/service/main/resource/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// Dao is resource dao.
type Dao struct {
	db *xsql.DB
	c  *conf.Config
	// redis
	redis  *redis.Pool
	expire int32
}

// New init mysql db
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     xsql.NewMySQL(c.DB.Ads),
		redis:  redis.NewPool(c.Redis.Ads.Config),
		expire: int32(time.Duration(c.Redis.Ads.Expire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}

// Ping check dao health.
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}
