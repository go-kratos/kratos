package dao

import (
	"context"

	"go-common/library/database/elastic"

	"go-common/app/admin/main/passport/conf"
)

// Dao dao
type Dao struct {
	c *conf.Config
	// elastic client
	EsCli *elastic.Elastic
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
		// elastic client
		EsCli: elastic.NewElastic(c.Elastic),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
