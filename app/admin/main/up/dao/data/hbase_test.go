package data

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatahbaseMd5Key(t *testing.T) {
	convey.Convey("hbaseMd5Key", t, func(ctx convey.C) {
		var (
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hbaseMd5Key(aid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataVideoQuitPoints(t *testing.T) {
	convey.Convey("VideoQuitPoints", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.VideoQuitPoints(c, cid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataArchiveStat(t *testing.T) {
	convey.Convey("ArchiveStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			stat, err := d.ArchiveStat(c, aid)
			ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(stat, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataArchiveArea(t *testing.T) {
	convey.Convey("ArchiveArea", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ArchiveArea(c, aid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataBaseUpStat(t *testing.T) {
	convey.Convey("BaseUpStat", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			date = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			stat, err := d.BaseUpStat(c, mid, date)
			ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(stat, convey.ShouldBeNil)
			})
		})
	})
}
