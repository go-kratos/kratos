package app

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppPgcCont(t *testing.T) {
	var (
		c     = context.Background()
		id    = int(0)
		limit = int(10)
	)
	convey.Convey("PgcCont", t, func(ctx convey.C) {
		res, err := d.PgcCont(c, id, limit)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestAppPgcContCount(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("PgcContCount", t, func(ctx convey.C) {
		upCnt, err := d.PgcContCount(c)
		ctx.Convey("Then err should be nil.upCnt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(upCnt, convey.ShouldNotBeNil)
		})
	})
}
