package http

import (
	"context"

	"go-common/app/admin/main/aegis/conf"
	"go-common/library/database/elastic"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c *conf.Config

	clientR, clientW *bm.Client
	es               *elastic.Elastic
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:       c,
		clientR: bm.NewClient(c.HTTPClient.Read),
		clientW: bm.NewClient(c.HTTPClient.Write),
		es: elastic.NewElastic(&elastic.Config{
			Host:       c.Host.Manager,
			HTTPClient: c.HTTPClient.Es,
		}),
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
