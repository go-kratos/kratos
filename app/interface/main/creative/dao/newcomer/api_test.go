package newcomer

import (
	"context"
	"go-common/library/ecode"
	"testing"

	"strings"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestNewcomerMall(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(27515308)
		couponID = "ef89b24b3951429b"
		uname    = "test"
	)
	convey.Convey("Mall", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.mallURI).Reply(200).JSON(`{"code":20062}`)
		err := d.Mall(c, mid, couponID, uname)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestNewcomerBCoin(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(27515406)
		money = int64(1)
		aid   = "217"
	)

	convey.Convey("BCoin", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.bPayURI).Reply(200).JSON(`{"code":20063}`)
		err := d.BCoin(c, mid, aid, money)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestNewcomerPendant(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(27515406)
		PID     = "4"
		expires = int64(1)
		err     error
	)
	convey.Convey("Pendant", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.pendantURI).Reply(200).JSON(`{"code":20064}`)
		err = d.Pendant(c, mid, PID, expires)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeNewcomerPendantAPIErr)
		})
	})
}

func TestNewcomerBigMemberCoupon(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(27515308)
		batchToken = "841764801720181030122848"
	)
	convey.Convey("BigMemberCoupon", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.bigMemberURI).Reply(200).JSON(`{"code":20069}`)
		err := d.BigMemberCoupon(c, mid, batchToken)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}

func TestNewcomerSendNotify(t *testing.T) {
	convey.Convey("SendNotify", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mids    = []int64{27515405}
			mc      = "1_17_4"
			title   = "creative"
			context = "sssss"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendNotify(c, mids, mc, title, context)
			ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
