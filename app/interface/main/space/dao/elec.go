package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_elecURI       = "/api/elec/info/query"
	_elecMonthRank = "1"
)

// ElecInfo .
func (d *Dao) ElecInfo(c context.Context, mid, paymid int64) (data *model.ElecInfo, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pay_mid", strconv.FormatInt(paymid, 10))
	params.Set("type", _elecMonthRank)
	var res struct {
		Code int             `json:"code"`
		Data *model.ElecInfo `json:"data"`
	}
	if err = d.httpR.Get(c, d.elecURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if res.Code == 500011 {
			return
		}
		err = errors.Wrap(ecode.Int(res.Code), d.elecURL+"?"+params.Encode())
		return
	}
	data = res.Data
	return
}
