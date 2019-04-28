package v2

import (
	"context"
	"time"

	"go-common/app/infra/config/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/orm"

	"github.com/jinzhu/gorm"
)

// Dao dao.
type Dao struct {
	// redis
	redis  *redis.Pool
	expire time.Duration
	// cache
	pathCache string
	//DB
	DB *gorm.DB
}

// New new a dao.
func New(c *conf.Config) *Dao {
	d := &Dao{
		// redis
		redis:  redis.NewPool(c.Redis),
		expire: time.Duration(c.PollTimeout),
		// cache
		pathCache: c.PathCache,
		// orm
		DB: orm.NewMySQL(c.ORM),
	}
	return d
}

// Ping ping is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	return d.DB.DB().PingContext(c)
}

// Close close resuouces.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}
