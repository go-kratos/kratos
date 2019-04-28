package dao

import (
	"context"

	"go-common/app/job/main/upload/conf"
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

}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
