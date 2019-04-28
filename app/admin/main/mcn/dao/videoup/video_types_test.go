package videoup

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVideouprefreshUpType(t *testing.T) {
	convey.Convey("refreshUpType", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.refreshUpType()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestVideoupGetTidName(t *testing.T) {
	convey.Convey("GetTidName", t, func(ctx convey.C) {
		var (
			tids = []int64{22}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tpNames := d.GetTidName(tids)
			ctx.Convey("Then tpNames should not be nil.", func(ctx convey.C) {
				ctx.So(tpNames, convey.ShouldNotBeNil)
			})
		})
	})
}
