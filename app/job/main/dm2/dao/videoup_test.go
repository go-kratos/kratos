package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoarchiveURI(t *testing.T) {
	convey.Convey("archiveURI", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.archiveURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoVideos(t *testing.T) {
	convey.Convey("Videos", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.Videos(c, aid)
		})
	})
}
