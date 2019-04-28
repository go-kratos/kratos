package notice

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNoticeTopic(t *testing.T) {
	convey.Convey("Topic", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.Topic(c, oid)
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

func TestNoticeActivity(t *testing.T) {
	convey.Convey("Activity", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.Activity(c, oid)
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

func TestNoticeActivitySub(t *testing.T) {
	convey.Convey("ActivitySub", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := d.ActivitySub(c, oid)
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
