package notice

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNoticeDrawyoo(t *testing.T) {
	convey.Convey("Drawyoo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			hid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.Drawyoo(c, hid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}
