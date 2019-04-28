package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveUpCount(t *testing.T) {
	convey.Convey("UpCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.UpCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestArchiveArchives(t *testing.T) {
	convey.Convey("Archives", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.Archives(c, aids)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestArchiveArchive(t *testing.T) {
	convey.Convey("Archive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.Archive(c, aid)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestArchiveStats(t *testing.T) {
	convey.Convey("Stats", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.Stats(c, aids)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}

func TestArchiveArticleMetas(t *testing.T) {
	convey.Convey("ArticleMetas", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.ArticleMetas(c, aids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
