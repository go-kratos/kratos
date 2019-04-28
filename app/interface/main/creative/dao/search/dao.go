package search

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/elastic"
	bm "go-common/library/net/http/blademaster"
)

// Dao is search dao.
type Dao struct {
	c *conf.Config
	// http client
	client *bm.Client
	// searchURI string
	// memberURI string
	es *elastic.Elastic
}

// New new a search dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// client
		client: bm.NewClient(c.HTTPClient.Slow),
		// uri
		// searchURI: c.Host.Search + _spaceURL,
		// memberURI: c.Host.Search + _memberURL,
		es: elastic.NewElastic(&elastic.Config{
			Host:       c.Host.MainSearch,
			HTTPClient: c.HTTPClient.Slow,
		}),
	}
	return d
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
