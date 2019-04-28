package bnj

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"go-common/library/ecode"

	"github.com/pkg/errors"
)

const _grantCouponURL = "/mall-marketing/coupon_code/create"

// GrantCoupon grant coupon to mid.
func (d *Dao) GrantCoupon(c context.Context, mid int64, couponID string) (err error) {
	var (
		bs  []byte
		req *http.Request
	)
	param := &struct {
		Mid      int64  `json:"mid"`
		CouponID string `json:"couponId"`
	}{
		Mid:      mid,
		CouponID: couponID,
	}
	if bs, err = json.Marshal(param); err != nil {
		return
	}
	if req, err = http.NewRequest(http.MethodPost, d.grantCouponURL, strings.NewReader(string(bs))); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.grantCouponURL+"msg:"+res.Msg)
	}
	return
}
