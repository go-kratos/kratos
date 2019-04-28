package dao

import (
	"context"

	"go-common/app/admin/main/cache/conf"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao dao.
type Dao struct {
	c      *conf.Config
	DB     *gorm.DB
	client *bm.Client
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		DB:     orm.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
	}
	return
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

// Close close connection of db , mc.
func (d *Dao) Close() {

}
