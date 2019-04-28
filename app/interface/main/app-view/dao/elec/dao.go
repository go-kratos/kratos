package elec

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/elec"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_elec          = "/api/elec/info/query"
	_elecTotal     = "/api/v2/rank/total/av/query"
	_elecMonthRank = "1"
)

// Dao is elec dao.
type Dao struct {
	client    *httpx.Client
	elec      string
	elecTotal string
}

// New elec dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    httpx.NewClient(c.HTTPClient),
		elec:      c.Host.Elec + _elec,
		elecTotal: c.Host.Elec + _elecTotal,
	}
	return
}

// TotalInfo mid+aid total elec info
func (d *Dao) TotalInfo(c context.Context, mid int64, aid int64) (data *elec.Info, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("av_no", strconv.FormatInt(aid, 10))
	var res struct {
		Code int        `json:"code"`
		Data *elec.Info `json:"data"`
	}
	if err = d.client.Get(c, d.elecTotal, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.elecTotal+"?"+params.Encode())
		return
	}
	data = res.Data
	return
}

// Info elec info
func (d *Dao) Info(c context.Context, mid, paymid int64) (data *elec.Info, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	if paymid > 0 {
		params.Set("pay_mid", strconv.FormatInt(paymid, 10))
	}
	params.Set("type", _elecMonthRank)
	var res struct {
		Code int        `json:"code"`
		Data *elec.Info `json:"data"`
	}
	if err = d.client.Get(c, d.elec, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if res.Code == 500011 {
			return
		}
		err = errors.Wrap(ecode.Int(res.Code), d.elec+"?"+params.Encode())
		return
	}
	data = res.Data
	return
}
