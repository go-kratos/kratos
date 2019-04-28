package card

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestListVideoArchive(t *testing.T) {
	var (
		c     = context.TODO()
		avids = []int64{31908629}
	)
	convey.Convey("Info", t, func(ctx convey.C) {
		videos, err := d.ListVideoArchive(c, avids)
		ctx.Convey("Then err should be nil.videos should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(videos, convey.ShouldNotBeNil)
		})
	})
}

func TestAvidVideoMap(t *testing.T) {
	var (
		c     = context.TODO()
		avids = []int64{31908629}
	)
	convey.Convey("Info", t, func(ctx convey.C) {
		avidVideoMap, err := d.AvidVideoMap(c, avids)
		ctx.Convey("Then err should be nil.midVideosMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(avidVideoMap, convey.ShouldNotBeNil)
		})
	})
}
