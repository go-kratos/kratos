package search

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchRecordPaginate(t *testing.T) {
	convey.Convey("RecordPaginate", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			types = []int64{}
			mid   = int64(0)
			stime = int64(0)
			etime = int64(0)
			order = ""
			sort  = ""
			pn    = int32(0)
			ps    = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			records, total, err := d.RecordPaginate(c, types, mid, stime, etime, order, sort, pn, ps)
			ctx.Convey("Then err should be nil.records,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(records, convey.ShouldBeNil)
			})
		})
	})
}
