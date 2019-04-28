package upcrmservice

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceGetDateStamp(t *testing.T) {
	convey.Convey("GetDateStamp", t, func(ctx convey.C) {
		var (
			timeStamp = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := GetDateStamp(timeStamp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
