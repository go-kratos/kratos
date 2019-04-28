package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestOverlordClusters(t *testing.T) {
	convey.Convey("get OverlordClusters", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			ocs, err := d.OverlordClusters(context.Background(), "", "main.common-arch.overlord")
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ocs, convey.ShouldNotBeNil)
			})
		})
	})
}
