package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertBgmWhiteList(t *testing.T) {
	convey.Convey("InsertBgmWhiteList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(111)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.InsertBgmWhiteList(c, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
