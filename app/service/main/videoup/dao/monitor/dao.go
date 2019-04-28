package monitor

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/service/main/videoup/conf"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

// Dao is message dao.
type Dao struct {
	c      *conf.Config
	client *xhttp.Client
	uri    string
}

// New new a message dao.
func New(c *conf.Config) (d *Dao) {
	//http://ops-mng.bilibili.co/api/sendsms&message=test&token=
	d = &Dao{
		c:      c,
		client: xhttp.NewClient(c.HTTPClient.Read),
		uri:    c.Host.Monitor + "/api/sendsms",
	}
	return
}

// Send send message to upper.
func (d *Dao) Send(c context.Context, msg string) (err error) {
	params := url.Values{}
	params.Set("phone", d.c.Monitor.Tels)
	params.Set("message", msg)
	params.Set("token", "f5a658b2-5926-4b71-96c3-7d3777b7d256")
	if err = d.client.Get(c, d.uri, "", params, nil); err != nil {
		log.Info("sms error(%v)", err)
		fmt.Printf("sms error(%v)", err)
	}
	return
}
