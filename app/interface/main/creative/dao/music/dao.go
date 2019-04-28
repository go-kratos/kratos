package music

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao is archive dao.
type Dao struct {
	// config
	c            *conf.Config
	db           *sql.DB
	client       *bm.Client
	audioListURL string
	es           *elastic.Elastic
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Archive),
		// client
		client:       bm.NewClient(c.HTTPClient.Slow),
		audioListURL: c.Host.API + _audioListURI,
		es: elastic.NewElastic(&elastic.Config{
			Host:       c.Host.MainSearch,
			HTTPClient: c.HTTPClient.Slow,
		}),
	}
	return
}

// Ping fn
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close fn
func (d *Dao) Close() (err error) {
	return d.db.Close()
}
