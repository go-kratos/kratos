package message

import (
	"context"
	"fmt"
	"go-common/app/service/main/assist/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/url"
	"strconv"
)

const (
	_messageURI = "/api/notify/send.user.notify.do"
)

// Dao is message dao.
type Dao struct {
	c      *conf.Config
	client *bm.Client
	uri    string
}

// New new a message dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: bm.NewClient(c.HTTPClient.Slow),
		uri:    c.Host.Message + _messageURI,
	}
	return d
}

// Send send message form uper to assistMid.
func (d *Dao) Send(c context.Context, mc, title, msg string, mid int64) (err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("source", "1")
	params.Set("data_type", "4")
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("context", msg)
	params.Set("mid_list", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.uri, "", params, &res); err != nil {
		log.Error("message url(%s) error(%v)", d.uri+"?"+params.Encode(), err)
		return
	}
	log.Info("SendSysNotify url: (%s)", d.uri+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("message url(%s) error(%v)", d.uri+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("message send failed")
	}
	return
}
