package hbase

import (
	"go-common/app/job/main/videoup-report/conf"
	"go-common/library/database/hbase.v2"
)

// Dao is redis dao.
type Dao struct {
	c     *conf.Config
	hbase *hbase.Client
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		hbase: hbase.NewClient(&c.Hbase.Config),
	}
	return d
}

// Close fn
func (d *Dao) Close() {
}
