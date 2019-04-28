package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAsoCleanCache(t *testing.T) {
	convey.Convey("AsoCleanCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			token   = ""
			session = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
			mid     = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AsoCleanCache(c, token, session, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
