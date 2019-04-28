package monitor

import (
	"context"
	"net/url"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/monitor"
	httpx "go-common/library/net/http/blademaster"
)

// Dao is message dao.
type Dao struct {
	c      *conf.Config
	client *httpx.Client
	uri    string
}

// New new a message dao.
func New(c *conf.Config) *Dao {
	//http://ops-mng.bilibili.co/api/sendsms&message=test&token=
	return &Dao{
		c:      c,
		client: httpx.NewClient(c.HTTPClient.Normal),
		uri:    c.Host.Monitor + "/api/sendsms",
	}
}

// Send send message to upper.
func (d *Dao) Send(c context.Context, msg string) (err error) {
	params := url.Values{}
	params.Set("phone", monitor.Tels)
	params.Set("message", msg)
	params.Set("token", "f5a658b2-5926-4b71-96c3-7d3777b7d256")
	d.client.Get(c, d.uri, "", params, nil)
	return
}
