package notice

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNoticeBangumi(t *testing.T) {
	convey.Convey("Bangumi", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, epid, err := d.Bangumi(c, oid)
			ctx.Convey("Then err should be nil.title,link,epid should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(epid, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}
