package shopping

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_couponURL = "/mall-marketing/coupon_code/create"
)

// Dao is shopping dao
type Dao struct {
	client    *httpx.Client
	couponURL string
}

// New shopping dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
		// url
		couponURL: c.Host.Mall + _couponURL,
	}
	return
}

// Coupon user vip
func (d *Dao) Coupon(c context.Context, couponID string, mid int64, uname string) (msg string, err error) {
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	data := map[string]interface{}{
		"couponId": couponID,
		"mid":      mid,
		"uname":    uname,
	}
	var (
		bytesData []byte
		req       *http.Request
	)
	if bytesData, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	if req, err = http.NewRequest("POST", d.couponURL, bytes.NewReader(bytesData)); err != nil {
		log.Error("http.NewRequest error(%v)", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("X-BACKEND-BILI-REAL-IP", "")
	log.Info("coupon vip mid(%d) couponID(%s)", mid, couponID)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("coupon vip url(%v) error(%v)", d.couponURL, err)
		return
	}
	msg = res.Msg
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("coupon vip url(%v) res code(%d)", d.couponURL, res.Code)
		return
	}
	return
}
