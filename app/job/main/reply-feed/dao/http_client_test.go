package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoIsOriginHot(t *testing.T) {
	convey.Convey("IsOriginHot", t, func(ctx convey.C) {
		var (
			oid  = int64(1)
			rpID = int64(1)
			tp   = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			isHot, err := d.IsOriginHot(context.Background(), oid, rpID, tp)
			ctx.Convey("Then err should be nil.isHot should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(isHot, convey.ShouldNotBeNil)
			})
		})
	})
}
