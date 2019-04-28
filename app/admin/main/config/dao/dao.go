package dao

import (
	"context"

	"go-common/app/admin/main/config/conf"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao dao.
type Dao struct {
	c      *conf.Config
	DB     *gorm.DB
	DBApm  *gorm.DB
	client *bm.Client
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		DB:     orm.NewMySQL(c.ORM),
		DBApm:  orm.NewMySQL(c.ORMApm),
		client: bm.NewClient(c.HTTPClient),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
	d.DBApm.LogMode(true)
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.DBApm != nil {
		d.DBApm.Close()
	}
}
