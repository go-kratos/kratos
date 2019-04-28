package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawVip(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110018671)
	)
	convey.Convey("Get vip-info from vip-rpc", t, func(ctx convey.C) {
		vip, err := d.RawVip(c, mid)
		ctx.Convey("Then err should be nil and vip should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vip, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawVips(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{110016841, 110018671}
	)
	convey.Convey("Batch get vip-infos from vip-rpc", t, func(ctx convey.C) {
		res, err := d.RawVips(c, mids)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
