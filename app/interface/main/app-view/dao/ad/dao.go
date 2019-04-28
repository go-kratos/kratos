package ad

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/ad"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_adURL          = "/bce/api/bce/wise"
	_monitorInfoURL = "/up-openapi/api/v1/av_monitor_info/%d"
)

// Dao dao.
type Dao struct {
	client         *bm.Client
	adURL          string
	monitorInfoURL string
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:         bm.NewClient(conf.Conf.HTTPAD),
		adURL:          c.Host.AD + _adURL,
		monitorInfoURL: c.Host.AD + _monitorInfoURL,
	}
	return
}

// Ad ad request.
func (d *Dao) Ad(c context.Context, mobiApp, device, buvid string, build int, mid, upperID, aid int64, rid int32, tids []int64, resource []int64, network, adExtra string) (advert *ad.Ad, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("device", device)
	params.Set("buvid", buvid)
	params.Set("resource", xstr.JoinInts(resource))
	params.Set("mobi_app", mobiApp)
	params.Set("ip", ip)
	params.Set("av_rid", strconv.FormatInt(int64(rid), 10))
	params.Set("av_tid", xstr.JoinInts(tids))
	params.Set("av_up_id", strconv.FormatInt(upperID, 10))
	if network != "" {
		params.Set("network", network)
	}
	if adExtra != "" {
		params.Set("ad_extra", adExtra)
	}
	var res struct {
		Code int    `json:"code"`
		Data *ad.Ad `json:"data"`
	}
	if err = d.client.Get(c, d.adURL, ip, params, &res); err != nil {
		return
	}
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		err = errors.Wrap(code, d.adURL+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		res.Data.ClientIP = ip
	}
	advert = res.Data
	return
}

// MonitorInfo ad aid monitor info
func (d *Dao) MonitorInfo(c context.Context, aid int64) (minfo json.RawMessage, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = d.client.RESTfulGet(c, d.monitorInfoURL, ip, nil, &res, aid); err != nil {
		return
	}
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		err = errors.Wrap(code, d.monitorInfoURL)
		return
	}
	minfo = res.Data
	return
}
