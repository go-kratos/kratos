package shop

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/shop"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const _info = "/api/merchants/shop/info"

type Dao struct {
	client *httpx.Client
	info   string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
		info:   c.Host.Show + _info,
	}
	return
}

func (d *Dao) Info(c context.Context, mid int64, mobiApp, device string, build int) (info *shop.Info, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("build", strconv.Itoa(build))
	var res struct {
		Code int        `json:"errno"`
		Data *shop.Info `json:"data"`
	}
	if err = d.client.Get(c, d.info, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if res.Code == 130000 {
			return
		}
		err = errors.Wrap(ecode.Int(res.Code), d.info+"?"+params.Encode())
		return
	}
	info = res.Data
	return
}
