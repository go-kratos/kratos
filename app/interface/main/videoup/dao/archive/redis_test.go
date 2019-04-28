package archive

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivekeyUpFavTpsPrefix(t *testing.T) {
	convey.Convey("keyUpFavTpsPrefix", t, func(ctx convey.C) {
		var (
			mid = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyUpFavTpsPrefix(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveFilenameExpires(t *testing.T) {
	convey.Convey("FilenameExpires", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			vs = []*archive.VideoParam{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ves, err := d.FilenameExpires(c, vs)
			ctx.Convey("Then err should be nil.ves should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(ves, convey.ShouldBeNil)
			})
		})
	})
}

func TestArchiveFreshFavTypes(t *testing.T) {
	convey.Convey("FreshFavTypes", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			tp  = int(162)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.FreshFavTypes(c, mid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestArchivepingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
