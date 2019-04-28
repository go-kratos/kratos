package dao

import (
	"go-common/app/admin/main/appstatic/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/orm"
	httpx "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao .
type Dao struct {
	DB     *gorm.DB
	c      *conf.Config
	client *httpx.Client
	redis  *redis.Pool
}

// New new a instance
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// db
		DB:     orm.NewMySQL(c.ORM),
		c:      c,
		client: httpx.NewClient(c.HTTPClient),
		redis:  redis.NewPool(c.Redis.Config),
	}
	d.DB.LogMode(true)
	return
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
