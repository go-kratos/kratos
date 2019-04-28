package search

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/ip"
)

// Dao is elec dao.
type Dao struct {
	c      *conf.Config
	client *httpx.Client
	key    string
	secret string
	notify string
}

// New pendant dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: httpx.NewClient(c.HTTPClient.Normal),
		notify: c.HTTPClient.SearchDomain + _notify,
	}
	return
}

const (
	_notify = "/sensitive/sync"
)

// Notify .
func (d *Dao) Notify(c context.Context, area []string) (err error) {
	ip := ip.InternalIP()
	params := url.Values{}
	params.Set("appkey", d.key)
	params.Set("appsecret", d.secret)
	params.Set("type", strings.Join(area, ","))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Get(c, d.notify, ip, params, &res); err != nil {
		log.Error("searchNotify url(%s) error(%v)", d.notify+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("searchNotify failed(%v)", res.Code)
		log.Error(" d.client.Get(%s) error(%v)", d.notify+"?"+params.Encode(), err)
		return
	}
	return
}
