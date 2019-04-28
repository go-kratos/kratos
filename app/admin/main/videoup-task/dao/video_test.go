package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetVID(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10098208)
		cid = int64(10108827)
	)
	convey.Convey("GetVID", t, func(ctx convey.C) {
		vid, err := d.GetVID(c, aid, cid)
		ctx.Convey("Then err should be nil.vid should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vid, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVideo(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110750)
		cid = int64(10134516)
	)
	convey.Convey("Video", t, func(ctx convey.C) {
		_, err := d.Video(c, aid, cid)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoVideoAttribute(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("VideoAttribute", t, func(ctx convey.C) {
		_, err := d.VideoAttribute(c, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoArcVideoByCID(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("ArcVideoByCID", t, func(ctx convey.C) {
		_, err := d.ArcVideoByCID(c, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoNewVideoByID(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("NewVideoByID", t, func(ctx convey.C) {
		_, err := d.NewVideoByID(c, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
