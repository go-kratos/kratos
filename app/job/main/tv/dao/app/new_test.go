package app

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppPromError(t *testing.T) {
	var (
		name = "test"
	)
	convey.Convey("PromError", t, func(ctx convey.C) {
		PromError(name)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestAppPromInfo(t *testing.T) {
	var (
		name = "test"
	)
	convey.Convey("PromInfo", t, func(ctx convey.C) {
		PromInfo(name)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestAppNumPce(t *testing.T) {
	var (
		count    = int(50)
		pagesize = int(10)
	)
	convey.Convey("NumPce", t, func(ctx convey.C) {
		numPce := NumPce(count, pagesize)
		ctx.Convey("Then numPce should not be nil.", func(ctx convey.C) {
			ctx.So(numPce, convey.ShouldNotBeNil)
		})
	})
}
