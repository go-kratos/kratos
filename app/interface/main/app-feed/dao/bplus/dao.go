package bplus

import (
	"context"
	"go-common/app/interface/main/app-card/model/bplus"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

const (
	_dynamicDetail = "/dynamic_detail/v0/dynamic/details"
)

// Dao is a dao.
type Dao struct {
	// http client
	client *httpx.Client
	// ad
	dynamicDetail string
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:        httpx.NewClient(c.HTTPClient),
		dynamicDetail: c.Host.DynamicCo + _dynamicDetail,
	}
	return
}

// DynamicDetail is.
func (d *Dao) DynamicDetail(c context.Context, ids ...int64) (picm map[int64]*bplus.Picture, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	for _, id := range ids {
		params.Add("dynamic_ids[]", strconv.FormatInt(id, 10))
	}
	var res struct {
		Code int `json:"code"`
		Data *struct {
			List []*bplus.Picture `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.dynamicDetail, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.dynamicDetail+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		picm = make(map[int64]*bplus.Picture, len(res.Data.List))
		for _, pic := range res.Data.List {
			picm[pic.DynamicID] = pic
		}
	}
	return
}
