package up

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpmcnSignKey(t *testing.T) {
	convey.Convey("mcnSignKey", t, func(ctx convey.C) {
		var (
			mcnMid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := mcnSignKey(mcnMid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpmcnUpperKey(t *testing.T) {
	convey.Convey("mcnUpperKey", t, func(ctx convey.C) {
		var (
			signID = int64(0)
			upMid  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := mcnUpperKey(signID, upMid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpDelMcnSignCache(t *testing.T) {
	convey.Convey("DelMcnSignCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcnMid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelMcnSignCache(c, mcnMid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpDelMcnUpperCache(t *testing.T) {
	convey.Convey("DelMcnUpperCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			upMid  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelMcnUpperCache(c, signID, upMid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
