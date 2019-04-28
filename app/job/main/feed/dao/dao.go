package dao

import (
	"context"
	"go-common/app/job/main/feed/conf"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

var (
	infosCount = prom.BusinessInfoCount
)

// Dao is feed job dao.
type Dao struct {
	c         *conf.Config
	smsClient *bm.Client
}

// New add a feed job dao.
func New(c *conf.Config) *Dao {
	return &Dao{
		c:         c,
		smsClient: bm.NewClient(c.HTTPClient),
	}
}

func (d *Dao) Ping(c context.Context) (err error) {
	return
}

func PromError(name string) {
	prom.BusinessErrCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}
