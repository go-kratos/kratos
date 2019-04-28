package manager

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerupSpecialCacheKey(t *testing.T) {
	convey.Convey("upSpecialCacheKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := upSpecialCacheKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
