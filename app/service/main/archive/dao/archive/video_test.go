package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivefirstCid(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
	)
	convey.Convey("firstCid", t, func(ctx convey.C) {
		cid, dimensions, err := d.firstCid(c, aid)
		ctx.Convey("Then err should be nil.cid,dimensions should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(dimensions, convey.ShouldNotBeNil)
			ctx.So(cid, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivevideos3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
	)
	convey.Convey("videos3", t, func(ctx convey.C) {
		_, err := d.videos3(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivevideosByAids3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{10097272}
	)
	convey.Convey("videosByAids3", t, func(ctx convey.C) {
		vs, err := d.videosByAids3(c, aids)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVideosByCids(t *testing.T) {
	var (
		c    = context.TODO()
		cids = []int64{10097272}
	)
	convey.Convey("VideosByCids", t, func(ctx convey.C) {
		vs, err := d.VideosByCids(c, cids)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivevideo3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
		cid = int64(10097272)
	)
	convey.Convey("video3", t, func(ctx convey.C) {
		_, err := d.video3(c, aid, cid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
