package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceNewUniqueCheck(t *testing.T) {
	convey.Convey("NewUniqueCheck", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := NewUniqueCheck()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceloadMcnUniqueCache(t *testing.T) {
	convey.Convey("loadMcnUniqueCache", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.loadMcnUniqueCache()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
