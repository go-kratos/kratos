package notice

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNoticeLiveSmallVideo(t *testing.T) {
	convey.Convey("LiveSmallVideo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.LiveSmallVideo(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNoticeLiveActivity(t *testing.T) {
	convey.Convey("LiveActivity", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.LiveActivity(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNoticeLiveNotice(t *testing.T) {
	convey.Convey("LiveNotice", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.LiveNotice(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNoticeLivePicture(t *testing.T) {
	convey.Convey("LivePicture", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.LivePicture(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}
