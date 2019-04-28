package elec

import (
	"context"

	"go-common/app/interface/main/videoup/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao is elec dao.
type Dao struct {
	c *conf.Config
	// http
	client *bm.Client
	// uri
	showURI     string
	arcOpenURL  string
	arcCloseURL string
}

// New new a elec dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// http client
		client: bm.NewClient(c.HTTPClient.Write),
		// uri
		showURI:     c.Host.Elec + _showURL,
		arcOpenURL:  c.Host.Elec + _arcOpenURI,
		arcCloseURL: c.Host.Elec + _arcCloseURI,
	}
	return d
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
