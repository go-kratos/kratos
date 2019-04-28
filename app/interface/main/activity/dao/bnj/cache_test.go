package bnj

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBnjtimeFinishKey(t *testing.T) {
	convey.Convey("timeFinishKey", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := timeFinishKey()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
