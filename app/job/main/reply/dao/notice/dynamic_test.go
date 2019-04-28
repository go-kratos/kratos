package notice

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNoticeDynamic(t *testing.T) {
	convey.Convey("Dynamic", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			content, link, err := d.Dynamic(c, oid)
			ctx.Convey("Then err should be nil.content,link should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(content, convey.ShouldNotBeNil)
			})
		})
	})
}
