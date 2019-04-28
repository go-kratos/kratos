package archive

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveBIZsByTime(t *testing.T) {
	convey.Convey("BIZsByTime", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = &time.Time{}
			end   = &time.Time{}
			tp    = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			bizs, err := d.BIZsByTime(c, start, end, tp)
			ctx.Convey("Then err should be nil.bizs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bizs, convey.ShouldNotBeNil)
			})
		})
	})
}
