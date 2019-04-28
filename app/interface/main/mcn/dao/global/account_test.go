package global

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGlobalGetInfo(t *testing.T) {
	convey.Convey("GetInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(417851)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := GetInfo(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetInfos(t *testing.T) {
	convey.Convey("GetInfos", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = []int64{417851}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := GetInfos(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetName(t *testing.T) {
	convey.Convey("GetName", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nickname := GetName(c, mid)
			ctx.Convey("Then nickname should not be nil.", func(ctx convey.C) {
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}
