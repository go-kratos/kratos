package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInfo3(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
	)
	convey.Convey("Info3", t, func(ctx convey.C) {
		res, err := d.Info3(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInfos(t *testing.T) {
	var (
		c      = context.Background()
		ipList = []string{"127.0.0.1"}
	)
	convey.Convey("Infos", t, func(ctx convey.C) {
		res, err := d.Infos(c, ipList)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCheckRealnameStatus(t *testing.T) {
	var (
		c         = context.Background()
		mid int64 = 1
	)
	convey.Convey("CheckRealnameStatus", t, func(ctx convey.C) {
		res, err := d.CheckRealnameStatus(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
