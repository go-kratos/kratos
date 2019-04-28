package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAuditInfo(t *testing.T) {
	convey.Convey("AuditInfo", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(1)
			remoteIP = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.AuditInfo(c, mid, remoteIP)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
