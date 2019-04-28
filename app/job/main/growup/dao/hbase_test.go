package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohbaseMd5Key(t *testing.T) {
	convey.Convey("hbaseMd5Key", t, func(ctx convey.C) {
		var (
			aid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hbaseMd5Key(aid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveStat(t *testing.T) {
	convey.Convey("ArchiveStat", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aid  = int64(2)
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			stat, err := d.ArchiveStat(c, aid, date)
			ctx.Convey("Then err should not be nil.stat should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(stat, convey.ShouldBeNil)
			})
		})
	})
}
