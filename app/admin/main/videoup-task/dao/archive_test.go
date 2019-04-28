package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoArchive(t *testing.T) {
	convey.Convey("Archive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Archive(c, aid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoArchiveParam(t *testing.T) {
	convey.Convey("ArchiveParam", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.ArchiveParam(c, aid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			})
		})
	})
}
