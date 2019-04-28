package dao

import (
	"context"

	"go-common/app/admin/main/sms/conf"
	"go-common/library/database/orm"
	xhttp "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao is the appeal database access object
type Dao struct {
	c          *conf.Config
	DB         *gorm.DB
	httpClient *xhttp.Client
}

// New will create a new appeal Dao instance
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		DB:         orm.NewMySQL(c.DB),
		httpClient: xhttp.NewClient(c.HTTPClient),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
}

// Close close dao.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	err = d.DB.DB().PingContext(c)
	return
}
