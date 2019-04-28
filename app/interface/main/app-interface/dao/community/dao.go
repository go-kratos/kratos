package community

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/community"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_comm = "/api/query.my.community.list.do"
)

// Dao is community dao
type Dao struct {
	client    *httpx.Client
	community string
}

// New initial community dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    httpx.NewClient(c.HTTPIm9),
		community: c.Host.Im9 + _comm,
	}
	return
}

// Community get community data from api.
func (d *Dao) Community(c context.Context, mid int64, ak, platform string, pn, ps int) (co []*community.Community, count int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("actionKey", "appkey")
	params.Set("data_type", "2")
	params.Set("access_key", ak)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("page_no", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	params.Set("platform", platform)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Count  int                    `json:"total_count"`
			Result []*community.Community `json:"result"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.community, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.community+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		co = res.Data.Result
		count = res.Data.Count
	}
	return
}
