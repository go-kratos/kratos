package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivelockKey(t *testing.T) {
	var (
		aid = int64(0)
	)
	convey.Convey("lockKey", t, func(ctx convey.C) {
		p1 := lockKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveGetVideoJam(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("GetVideoJam", t, func(ctx convey.C) {
		seconds, err := d.GetVideoJam(c)
		ctx.Convey("Then err should be nil.seconds should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(seconds, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveSetNXLock(t *testing.T) {
	var (
		c     = context.Background()
		aid   = int64(0)
		times = int64(0)
	)
	convey.Convey("SetNXLock", t, func(ctx convey.C) {
		res, err := d.SetNXLock(c, aid, times)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveDelLock(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(0)
	)
	convey.Convey("DelLock", t, func(ctx convey.C) {
		err := d.DelLock(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
