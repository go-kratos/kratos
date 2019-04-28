package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetPGCInfo(t *testing.T) {
	convey.Convey("GetPGCInfo", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			oid = int32(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.GetPGCInfo(c, oid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
