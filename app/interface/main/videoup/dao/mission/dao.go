package mission

import (
	"context"

	"go-common/app/interface/main/videoup/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao is mission dao.
type Dao struct {
	c                  *conf.Config
	httpR              *bm.Client
	missAllURL         string
	actOnlineByTypeURL string
}

// New new a mission dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// client
		httpR: bm.NewClient(c.HTTPClient.Read),
		// uri
		missAllURL:         c.Host.WWW + _msAllURL,
		actOnlineByTypeURL: c.Host.WWW + _actOnlineByTypeURI,
	}
	return d
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
