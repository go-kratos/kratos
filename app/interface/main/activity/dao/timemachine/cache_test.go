package timemachine

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTimemachinetimemachineKey(t *testing.T) {
	convey.Convey("timemachineKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := timemachineKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
