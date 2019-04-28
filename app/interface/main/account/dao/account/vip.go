package account

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_vipPointBalance = "/internal/v1/point/"
)

// VipPointBalance vip point
func (d *Dao) VipPointBalance(c context.Context, mid int64, ip string) (pointBalance int64, err error) {
	params := url.Values{}
	var res struct {
		Code int `json:"code"`
		Data struct {
			PointBalance int64 `json:"pointBalance"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.vipHost+_vipPointBalance+strconv.FormatInt(mid, 10), ip, params, &res); err != nil {
		log.Error("VipPointBalance() failed url(%v), err(%v)", d.vipHost+_vipPointBalance, err)
		return
	}
	if res.Code != 0 {
		log.Error("VipPointBalance() url(%v), res(%v)", d.vipHost+_vipPointBalance, res)
		err = ecode.Int(res.Code)
		return
	}
	pointBalance = res.Data.PointBalance
	return
}
