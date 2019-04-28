package upcrm

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmQueryPlayInfo(t *testing.T) {
	convey.Convey("QueryPlayInfo", t, func(ctx convey.C) {
		var (
			mid      = int64(0)
			busiType = []int{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryPlayInfo(mid, busiType)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryPlayInfoBatch(t *testing.T) {
	convey.Convey("QueryPlayInfoBatch", t, func(ctx convey.C) {
		var (
			mid      = []int64{}
			busiType = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryPlayInfoBatch(mid, busiType)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
