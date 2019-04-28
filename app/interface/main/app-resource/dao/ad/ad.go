package ad

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/splash"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_splashListURL = "/bce/api/splash/list"
	_splashShowURL = "/bce/api/splash/show"
)

// Dao is advertising dao.
type Dao struct {
	client        *httpx.Client
	splashListURL string
	splashShowURL string
}

// New advertising dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:        httpx.NewClient(conf.Conf.HTTPClient),
		splashListURL: c.Host.Ad + _splashListURL,
		splashShowURL: c.Host.Ad + _splashShowURL,
	}
	return
}

// SplashList ad splash list
func (d *Dao) SplashList(c context.Context, mobiApp, device, buvid, birth, adExtra string, height, width, build int, mid int64) (res []*splash.List, config *splash.CmConfig, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("build", strconv.Itoa(build))
	params.Set("buvid", buvid)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("ip", ip)
	params.Set("height", strconv.Itoa(height))
	params.Set("width", strconv.Itoa(width))
	params.Set("mid", strconv.FormatInt(mid, 10))
	if birth != "" {
		params.Set("birth", birth)
	}
	if adExtra != "" {
		params.Set("ad_extra", adExtra)
	}
	var data struct {
		Code int `json:"code"`
		*splash.CmConfig
		RequestID string         `json:"request_id"`
		Data      []*splash.List `json:"data"`
	}
	if err = d.client.Get(c, d.splashListURL, ip, params, &data); err != nil {
		log.Error("cpm splash url(%s) error(%v)", d.splashListURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&data)
	log.Info("cpm splash list url(%s) response(%s)", d.splashListURL+"?"+params.Encode(), b)
	if data.Code != 0 {
		err = ecode.Int(data.Code)
		log.Error("cpm splash url(%s) code(%d)", d.splashListURL+"?"+params.Encode(), data.Code)
		return
	}
	for _, t := range data.Data {
		s := &splash.List{}
		*s = *t
		s.RequestID = data.RequestID
		s.ClientIP = ip
		s.IsAdLoc = true
		res = append(res, s)
	}
	config = data.CmConfig
	return
}

// SplashShow ad splash show
func (d *Dao) SplashShow(c context.Context, mobiApp, device, buvid, birth, adExtra string, height, width, build int, mid int64) (res []*splash.Show, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("build", strconv.Itoa(build))
	params.Set("buvid", buvid)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("ip", ip)
	params.Set("height", strconv.Itoa(height))
	params.Set("width", strconv.Itoa(width))
	params.Set("mid", strconv.FormatInt(mid, 10))
	if birth != "" {
		params.Set("birth", birth)
	}
	if adExtra != "" {
		params.Set("ad_extra", adExtra)
	}
	var data struct {
		Code int            `json:"code"`
		Data []*splash.Show `json:"data"`
	}
	if err = d.client.Get(c, d.splashShowURL, ip, params, &data); err != nil {
		log.Error("cpm splash url(%s) error(%v)", d.splashShowURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&data)
	log.Info("cpm splash show url(%s) response(%s)", d.splashShowURL+"?"+params.Encode(), b)
	if data.Code != 0 {
		err = ecode.Int(data.Code)
		log.Error("cpm splash url(%s) code(%d)", d.splashShowURL+"?"+params.Encode(), data.Code)
		return
	}
	res = data.Data
	return
}
