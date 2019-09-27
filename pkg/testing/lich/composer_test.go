package lich

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestComposer(t *testing.T) {
	convey.Convey("Composer testing....", t, func(convCtx convey.C) {
		convCtx.Convey("When Setup everything goes positive", func(convCtx convey.C) {
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(Setup(), convey.ShouldBeNil)
			})
		})

		convCtx.Convey("When UnSetup everything goes positive", func(convCtx convey.C) {
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(Teardown(), convey.ShouldBeNil)
			})
		})
	})
}
