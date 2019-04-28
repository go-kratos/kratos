package music

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMusicCategorys(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{10, 11, 12}
	)
	convey.Convey("Categorys", t, func(ctx convey.C) {
		res, resMap, err := d.Categorys(c, ids)
		ctx.Convey("Then err should be nil.res,resMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resMap, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestMusicMCategorys(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("MCategorys", t, func(ctx convey.C) {
		res, err := d.MCategorys(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestMusicMusic(t *testing.T) {
	var (
		c    = context.TODO()
		sids = []int64{1, 2, 3, 4, 5, 6}
	)
	convey.Convey("Music", t, func(ctx convey.C) {
		res, err := d.Music(c, sids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
