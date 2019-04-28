package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveArchive(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110560)
		ip  = "127.0.0.1"
	)
	convey.Convey("Archive", t, func(ctx convey.C) {
		a, err := d.Archive(c, aid, ip)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveArchives(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{}
		ip   = ""
	)
	convey.Convey("Archives", t, func(ctx convey.C) {
		a, err := d.Archives(c, aids, ip)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveStats(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{}
		ip   = ""
	)
	convey.Convey("Stats", t, func(ctx convey.C) {
		a, err := d.Stats(c, aids, ip)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpCount(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(888952460)
	)
	convey.Convey("UpCount", t, func(ctx convey.C) {
		count, err := d.UpCount(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVideo(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110816)
		cid = int64(10134702)
		ip  = "127.0.0.1"
	)
	convey.Convey("Video", t, func(ctx convey.C) {
		v, err := d.Video(c, aid, cid, ip)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(v, convey.ShouldNotBeNil)
		})
	})
}
