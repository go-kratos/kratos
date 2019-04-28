package academy

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAcademyArchive(t *testing.T) {
	convey.Convey("Archive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			bs  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			a, err := d.Archive(c, oid, bs)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyArchiveCount(t *testing.T) {
	convey.Convey("ArchiveCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{}
			bs   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.ArchiveCount(c, tids, bs)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademySearchArchive(t *testing.T) {
	convey.Convey("SearchArchive", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tidsMap map[int][]int64
			bs      = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SearchArchive(c, tidsMap, bs)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyArchiveTagsByOids(t *testing.T) {
	convey.Convey("ArchiveTagsByOids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ArchiveTagsByOids(c, oids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
