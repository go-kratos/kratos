package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLatestNotice(t *testing.T) {
	convey.Convey("LatestNotice", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			platform = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO notice(status, platform, title, type) VALUES(1, 3, 'test', 0)")
			n, err := d.LatestNotice(c, platform)
			ctx.Convey("Then err should be nil.n should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(n, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoNotices(t *testing.T) {
	convey.Convey("Notices", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			typ      = ""
			platform = int(3)
			offset   = int(0)
			limit    = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO notice(status, platform, title, type) VALUES(1, 3, 'test', 0)")
			notices, err := d.Notices(c, typ, platform, offset, limit)
			ctx.Convey("Then err should be nil.notices should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(notices, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoNoticeCount(t *testing.T) {
	convey.Convey("NoticeCount", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			typ      = ""
			platform = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO notice(status, platform, title, type) VALUES(1, 3, 'test', 0)")
			count, err := d.NoticeCount(c, typ, platform)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
