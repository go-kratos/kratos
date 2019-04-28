package search

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_upper = "/main/recommend"
)

type Dao struct {
	client *httpx.Client
	upper  string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPSearch),
		upper:  c.Host.Search + _upper,
	}
	return
}

func (d *Dao) Follow(c context.Context, platform, mobiApp, device, buvid string, build int, mid, vmid int64) (ups []*search.Upper, trackID string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("clientip", ip)
	params.Set("build", strconv.Itoa(build))
	params.Set("buvid", buvid)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("context_id", strconv.FormatInt(vmid, 10))
	params.Set("rec_type", "up_rec")
	params.Set("pagesize", "20")
	params.Set("service_area", "play_suggest")
	var res struct {
		Code    int             `json:"code"`
		TrackID string          `json:"trackid"`
		Msg     string          `json:"msg"`
		Data    []*search.Upper `json:"data"`
	}
	if err = d.client.Get(c, d.upper, ip, params, &res); err != nil {
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("search list url(%s) response(%s)", d.upper+"?"+params.Encode(), b)
	code := ecode.Int(res.Code)
	if !code.Equal(ecode.OK) {
		err = errors.Wrap(code, d.upper+"?"+params.Encode())
		return
	}
	ups = res.Data
	trackID = res.TrackID
	return
}
