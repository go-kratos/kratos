package appeal

import (
	"go-common/app/interface/main/creative/conf"
	httpx "go-common/library/net/http/blademaster"
)

// Dao  define
type Dao struct {
	// config
	c *conf.Config
	// http client
	client *httpx.Client
	// appeal list
	list           string
	detail         string
	addappeal      string
	addreply       string
	appealstar     string
	appealstate    string
	appealStarInfo string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// http client
		client:         httpx.NewClient(c.HTTPClient.Normal),
		list:           c.Host.API + _list,
		detail:         c.Host.API + _detail,
		addappeal:      c.Host.API + _addappeal,
		addreply:       c.Host.API + _addreply,
		appealstar:     c.Host.API + _appealstar,
		appealstate:    c.Host.API + _appealstate,
		appealStarInfo: c.Host.API + _appealStarInfo,
	}
	return
}
