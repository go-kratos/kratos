package pay

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/elastic"
	bm "go-common/library/net/http/blademaster"
)

// Dao str
type Dao struct {
	c      *conf.Config
	client *bm.Client
	assURI string
	es     *elastic.Elastic
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: bm.NewClient(c.HTTPClient.Normal),
		assURI: c.Host.API + _assURI,
		es: elastic.NewElastic(&elastic.Config{
			Host:       c.Host.MainSearch,
			HTTPClient: c.HTTPClient.Slow,
		}),
	}
	return d
}

// Ping fn
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
