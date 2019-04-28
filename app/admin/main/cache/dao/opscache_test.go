package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestOpsMemcaches(t *testing.T) {
	convey.Convey("get OpsMemcaches", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			mcs, err := d.OpsMemcaches(context.Background())
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mcs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestOpsRediss(t *testing.T) {
	convey.Convey("get OpsRediss", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			mcs, err := d.OpsRediss(context.Background())
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mcs, convey.ShouldNotBeNil)
			})
		})
	})
}
