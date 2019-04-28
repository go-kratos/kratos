package dao

import (
	"context"

	"go-common/app/service/main/open/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/orm"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// Dao .
type Dao struct {
	// db
	DB *gorm.DB
	//memcache
	mc *memcache.Pool
}

// New new a instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// db
		DB: orm.NewMySQL(c.ORM),
		//memcache
		mc: memcache.NewPool(c.Memcache),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
}

// Ping check connection of db , mc .
func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		if err = d.DB.DB().PingContext(c); err != nil {
			log.Error("d.PingContext error (%v)", err)
		}
	}
	if err = d.pingMC(c); err != nil {
		return
	}
	return
}

// Close close connection of db.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
