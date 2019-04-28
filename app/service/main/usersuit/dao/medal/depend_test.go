package medal

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMedalSendMsg(t *testing.T) {
	convey.Convey("SendMsg", t, func(ctx convey.C) {
		var (
			mid     = int64(88889017)
			title   = ""
			context = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendMsg(c, mid, title, context)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalGetWearedfansMedal(t *testing.T) {
	convey.Convey("GetWearedfansMedal", t, func(ctx convey.C) {
		var (
			mid    = int64(88889017)
			source = int8(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			isLove, err := d.GetWearedfansMedal(c, mid, source)
			ctx.Convey("Then err should be nil.isLove should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(isLove, convey.ShouldNotBeNil)
			})
		})
	})
}
