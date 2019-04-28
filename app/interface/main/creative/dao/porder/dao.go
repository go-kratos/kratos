package porder

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/porder"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"net/url"
)

const (
	_porderConfig = "/videoup/porder/config/list"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	porderConfigURL string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		client:          httpx.NewClient(c.HTTPClient.Normal),
		porderConfigURL: c.Host.Videoup + _porderConfig,
	}
	return
}

// ListConfig fn
func (d *Dao) ListConfig(c context.Context) (cfgs []*porder.Config, err error) {
	params := url.Values{}
	var res struct {
		Code int              `json:"code"`
		Cfgs []*porder.Config `json:"data"`
	}
	if err = d.client.Get(c, d.porderConfigURL, "", params, &res); err != nil {
		log.Error("ListConfig url(%s) response(%+v) error(%v)", d.porderConfigURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	log.Info("ListConfig url(%s)", d.porderConfigURL+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("ListConfig url(%s) res(%v)", d.porderConfigURL, res)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	cfgs = res.Cfgs
	return
}
