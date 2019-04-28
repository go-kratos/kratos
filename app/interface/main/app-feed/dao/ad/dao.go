package ad

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_bce = "/bce/api/bce/wise"
)

// Dao is ad dao.
type Dao struct {
	// http client
	client *httpx.Client
	// ad
	bce string
}

// New new a ad dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client: httpx.NewClient(c.HTTPAd),
		bce:    c.Host.Ad + _bce,
	}
	return
}

func (d *Dao) Ad(c context.Context, mid int64, build int, buvid string, resource []int64, country, province, city, network, mobiApp, device, openEvent, adExtra string, style int, now time.Time) (advert *cm.Ad, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid)
	params.Set("resource", xstr.JoinInts(resource))
	params.Set("ip", ip)
	params.Set("country", country)
	params.Set("province", province)
	params.Set("city", city)
	params.Set("network", network)
	params.Set("build", strconv.Itoa(build))
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("open_event", openEvent)
	params.Set("ad_extra", adExtra)
	// 老接口做兼容
	if style != 0 {
		if style != 2 {
			style = 1
		}
		params.Set("style", strconv.Itoa(style))
	}
	var res struct {
		Code int    `json:"code"`
		Data *cm.Ad `json:"data"`
	}
	if err = d.client.Get(c, d.bce, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.bce+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		res.Data.ClientIP = ip
	}
	advert = res.Data
	return
}
