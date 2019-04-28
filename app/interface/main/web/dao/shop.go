package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_shopURI    = "/mall-shop/merchant/enter/service/shop/get"
	_shopTypePc = "2"
)

// ShopInfo get shop info data.
func (d *Dao) ShopInfo(c context.Context, mid int64) (data *model.ShopInfo, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("type", _shopTypePc)
	var res struct {
		Code int             `json:"code"`
		Data *model.ShopInfo `json:"data"`
	}
	if err = d.httpR.Get(c, d.shopURL, ip, params, &res); err != nil {
		err = errors.Wrapf(err, "ShopInfo(%s) mid(%d)", d.shopURL+params.Encode(), mid)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}
