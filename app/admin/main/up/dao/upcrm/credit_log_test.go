package upcrm

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmGetCreditLog(t *testing.T) {
	convey.Convey("GetCreditLog", t, func(ctx convey.C) {
		var (
			mid   = int64(1)
			limit = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetCreditLog(mid, limit)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldBeNil)
			})
		})
	})
}
