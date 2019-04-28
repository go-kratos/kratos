package pendant

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestPendantPayBcoin(t *testing.T) {
	convey.Convey("PayBcoin", t, func(ctx convey.C) {
		var (
			params   = url.Values{}
			ip       = ""
			_subject = "头像挂件"
		)
		params.Set("mid", "109228")
		params.Set("out_trade_no", "2016050614625209018624230766")
		params.Set("money", strconv.FormatFloat(2, 'f', 2, 64))
		params.Set("subject", _subject)
		params.Set("remark", fmt.Sprintf(_subject+" - %s（%s个月）", strconv.FormatInt(4, 10), strconv.FormatInt(1234, 10)))
		params.Set("merchant_id", d.c.PayInfo.MerchantID)
		params.Set("merchant_product_id", d.c.PayInfo.MerchantProductID)
		params.Set("platform_type", "3")
		params.Set("iap_pay_type", "0")
		params.Set("notify_url", d.c.PayInfo.CallBackURL)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("Post", d.payURL).Reply(0).JSON(`{"code":0}`)
			orderNo, casherURL, err := d.PayBcoin(c, params, ip)
			ctx.Convey("Then err should be nil.orderNo,casherURL should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(casherURL, convey.ShouldNotBeNil)
				ctx.So(orderNo, convey.ShouldNotBeNil)
			})
		})
	})
}
