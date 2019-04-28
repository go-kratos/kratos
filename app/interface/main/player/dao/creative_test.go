package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoViewPoints(t *testing.T) {
	convey.Convey("ViewPoints", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(10110670)
			cid = int64(10134319)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			points, err := d.ViewPoints(c, aid, cid)
			ctx.Convey("Then err should be nil.points should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(points, convey.ShouldNotBeNil)
				ctx.Printf("%+v", points)
			})
		})
	})
}
