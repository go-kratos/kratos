package monitor

import (
	"go-common/app/job/main/app/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	c      *conf.Config
	client *bm.Client
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: bm.NewClient(c.HTTPClient),
	}
	return d
}
