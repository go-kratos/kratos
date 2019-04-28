package dao

import (
	"context"

	"go-common/app/service/bbq/video-image/conf"
)

// Dao dao
type Dao struct {
	c *conf.Config
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return nil
}
