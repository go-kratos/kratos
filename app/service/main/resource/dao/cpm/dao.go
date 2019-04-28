package cpm

import (
	"go-common/app/service/main/resource/conf"
	httpx "go-common/library/net/http/blademaster"
)

// Dao define db struct
type Dao struct {
	c *conf.Config
	// cpt
	httpClient *httpx.Client
	cpmPCURL   string
	cpmAppURL  string
}

const (
	_cpmPCURL  = "/bce/api/bce/pc"
	_cpmAppURL = "/bce/api/bce/wise"
)

// New init mysql db
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		httpClient: httpx.NewClient(c.HTTPClient),
		cpmPCURL:   c.Host.Ad + _cpmPCURL,
		cpmAppURL:  c.Host.Ad + _cpmAppURL,
	}
	return
}
