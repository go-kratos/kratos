package pendant

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
)

// PayBcoin pay coin
func (d *Dao) PayBcoin(c context.Context, params url.Values, ip string) (orderNo, casherURL string, err error) {
	var res struct {
		Code      int    `json:"code"`
		Ts        string `json:"ts"`
		OrderNum  string `json:"order_no"`
		CasherURL string `json:"cashier_url"`
	}
	if err = d.client.Post(c, d.payURL, ip, params, &res); err != nil {
		log.Error("dao.client.Post(%s) error(%v)", d.payURL, err)
		return
	}
	if res.Code != 0 {
		log.Error("dao.client.Post(%s) error(%v)", d.payURL, res)
		err = ecode.Int(res.Code)
		return
	}
	orderNo = res.OrderNum
	casherURL = res.CasherURL
	return
}
