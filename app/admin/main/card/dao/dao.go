package dao

import (
	"context"

	"go-common/app/admin/main/card/conf"
	"go-common/library/database/orm"

	"github.com/jinzhu/gorm"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	DB *gorm.DB
}

// New init mysql db
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// db
		DB: orm.NewMySQL(c.ORM),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
}

// Close close the resource.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}
