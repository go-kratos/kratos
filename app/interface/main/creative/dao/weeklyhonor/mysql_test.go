package weeklyhonor

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWeeklyhonorpingMySQL(t *testing.T) {
	convey.Convey("pingMySQL", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMySQL(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWeeklyhonorHonorLogs(t *testing.T) {
	convey.Convey("HonorLogs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			hls, err := d.HonorLogs(c, mid)
			ctx.Convey("Then err should be nil.hls should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(hls, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorUpsertCount(t *testing.T) {
	convey.Convey("UpsertCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			hid = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpsertCount(c, mid, hid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWeeklyhonorUpsertClickCount(t *testing.T) {
	convey.Convey("UpsertClickCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(500)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpsertClickCount(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
