package kfc

import (
	"go-common/app/job/main/activity/conf"
	"go-common/library/net/http/blademaster"
)

// Dao .
type Dao struct {
	httpClient *blademaster.Client
	kfcDelURL  string
}

// New init
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		httpClient: blademaster.NewClient(c.HTTPClient),
		kfcDelURL:  c.Host.APICo + _kfcDelURI,
	}
	return
}
