package vip

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_vipActive = "/internal/v1/notice/active"
)

// Dao is vip dao.
type Dao struct {
	client       *httpx.Client
	vipActiveURL string
}

// New vip dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPWrite),
		// api
		vipActiveURL: c.Host.VIP + _vipActive,
	}
	return
}

// VIPActive get vip active info.
func (d *Dao) VIPActive(c context.Context, subID int) (msg string, err error) {
	params := url.Values{}
	params.Set("subId", strconv.Itoa(subID))
	var res struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err = d.client.Get(c, d.vipActiveURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.vipActiveURL+"?"+params.Encode())
		return
	}
	msg = res.Data
	return
}
