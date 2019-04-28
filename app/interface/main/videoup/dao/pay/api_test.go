package pay

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPayAssReg(t *testing.T) {
	convey.Convey("AssReg", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			aid = int64(10110826)
			bp  = int(5)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AssReg(c, mid, aid, bp, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPayAss(t *testing.T) {
	convey.Convey("Ass", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(10110826)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			assert, registed, err := d.Ass(c, aid, ip)
			ctx.Convey("Then err should be nil.assert,registed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(registed, convey.ShouldNotBeNil)
				ctx.So(assert, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPayUserAcceptProtocol(t *testing.T) {
	convey.Convey("UserAcceptProtocol", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			protocolID = "iamhashstring"
			mid        = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			accept, err := d.UserAcceptProtocol(c, protocolID, mid)
			ctx.Convey("Then err should be nil.accept should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(accept, convey.ShouldNotBeNil)
			})
		})
	})
}
