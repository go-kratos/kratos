package dao

import (
	"context"

	"go-common/app/interface/main/push/conf"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	c           *conf.Config
	reportPub   *databus.Databus
	callbackPub *databus.Databus
}

// New creates a push-service DAO instance.
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:           c,
		reportPub:   databus.New(c.ReportPub),
		callbackPub: databus.New(c.CallbackPub),
	}
	return d
}

// PromError prom error
func PromError(name string) {
	prom.BusinessErrCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	prom.BusinessInfoCount.Incr(name)
}

// Close dao.
func (d *Dao) Close() {}

// Ping check connection status.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
