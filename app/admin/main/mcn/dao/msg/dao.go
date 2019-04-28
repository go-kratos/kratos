package msg

import (
	"go-common/app/admin/main/mcn/conf"
	bm "go-common/library/net/http/blademaster"
)

const (
	_msgURL = "/api/notify/send.user.notify.do"
)

// Dao .
type Dao struct {
	c      *conf.Config
	client *bm.Client
	msgURL string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// http client
		client: bm.NewClient(c.HTTPClient),
		msgURL: c.Host.Msg + _msgURL,
	}
	return
}
