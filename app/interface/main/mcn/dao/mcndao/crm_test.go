package mcndao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaoGetActiveTid(t *testing.T) {
	convey.Convey("GetActiveTid", t, func(ctx convey.C) {
		var (
			mids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetActiveTid(mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
