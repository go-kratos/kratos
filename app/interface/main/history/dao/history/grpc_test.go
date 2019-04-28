package history

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestHistoryBangumis(t *testing.T) {
	convey.Convey("Bangumis", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			epid = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Bangumis(c, mid, epid)
		})
	})
}

func TestHistoryBangumisByAids(t *testing.T) {
	convey.Convey("BangumisByAids", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(14771787)
			aids   = []int64{11, 33}
			realIP = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.BangumisByAids(c, mid, aids, realIP)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
