package ad

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/ad"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_adURL = "/bce/api/bce/wise"
)

// Dao is advertising dao.
type Dao struct {
	client *httpx.Client
	adURL  string
}

// New advertising dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(conf.Conf.HTTPClient),
		adURL:  c.Host.Ad + _adURL,
	}
	return
}

// ADRequest Banners
func (d *Dao) ADRequest(c context.Context, mid int64, build int, buvid, resource, ip, country, province, city, network, mobiApp, device string, isPost bool) (adr *ad.ADRequest, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("buvid", buvid)
	params.Set("resource", resource)
	params.Set("mobi_app", mobiApp)
	params.Set("ip", ip)
	if device != "" {
		params.Set("device", device)
	}
	if country != "" {
		params.Set("country", country)
	}
	if province != "" {
		params.Set("province", province)
	}
	if city != "" {
		params.Set("city", city)
	}
	if network != "" {
		params.Set("network", network)
	}
	var res struct {
		Code int           `json:"code"`
		Data *ad.ADRequest `json:"data"`
	}
	if isPost {
		if err = d.client.Post(c, d.adURL, ip, params, &res); err != nil {
			log.Error("ad url(%s) error(%v)", d.adURL+"?"+params.Encode(), err)
			return
		}
	} else {
		if err = d.client.Get(c, d.adURL, ip, params, &res); err != nil {
			log.Error("ad url(%s) error(%v)", d.adURL+"?"+params.Encode(), err)
			return
		}
	}
	if res.Code != 0 {
		err = fmt.Errorf("ad api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d) or res.data(%v)", d.adURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	adr = res.Data
	return
}
