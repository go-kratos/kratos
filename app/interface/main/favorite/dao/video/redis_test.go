package video

import (
	"context"
	"testing"

	"go-common/app/service/main/favorite/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestVideocoversKey(t *testing.T) {
	convey.Convey("coversKey", t, func(ctx convey.C) {
		var (
			mid = int64(1)
			fid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := coversKey(mid, fid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func Test_pingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.pingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestVideoSetNewCoverCache(t *testing.T) {
	convey.Convey("SetNewCoverCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(1)
			fid    = int64(1)
			covers = []*model.Cover{{
				Aid:  123,
				Pic:  "123",
				Type: 2,
			},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetNewCoverCache(c, mid, fid, covers)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestVideoNewCoversCache(t *testing.T) {
	convey.Convey("NewCoversCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(1)
			fids = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, _, err := d.NewCoversCache(c, mid, fids)
			ctx.Convey("Then err should be nil.fcvs,mis should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDelCoverCache(t *testing.T) {
	convey.Convey("DelCoverCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			fid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCoverCache(c, mid, fid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
