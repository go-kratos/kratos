package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSearchVideo(t *testing.T) {
	var (
		c     = context.Background()
		pn    = int(1)
		ps    = int(10)
		query = "aaaaa"
	)
	convey.Convey("SearchVideo", t, func(ctx convey.C) {
		res, count, err := d.SearchVideo(c, pn, ps, query)
		ctx.Convey("Then err should be nil.res,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
