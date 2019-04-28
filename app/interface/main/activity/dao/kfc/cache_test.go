package kfc

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestKfckfcKey(t *testing.T) {
	convey.Convey("kfcKey", t, func(convCtx convey.C) {
		var (
			id = int64(3)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := kfcKey(id)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestKfckfcCodeKey(t *testing.T) {
	convey.Convey("kfcCodeKey", t, func(convCtx convey.C) {
		var (
			code = "201812041201"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := kfcCodeKey(code)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
