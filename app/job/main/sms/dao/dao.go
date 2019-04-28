package dao

import (
	"context"

	"go-common/app/job/main/sms/conf"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao struct info of Dao.
type Dao struct {
	c          *conf.Config
	httpClient *bm.Client
}

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		httpClient: bm.NewClient(c.HTTPClient),
	}
	return
}

// PromError prometheus error count.
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo prometheus info count.
func PromInfo(name string) {
	infosCount.Incr(name)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
