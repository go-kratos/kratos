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
			mid = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.HonorLogs(c, mid)
			ctx.Convey("Then err should be nil.hls should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWeeklyhonorLatestHonorLogs(t *testing.T) {
	convey.Convey("LatestHonorLogs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.LatestHonorLogs(c, mids)
			ctx.Convey("Then err should be nil.hls should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
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
			affected, err := d.UpsertCount(c, mid, hid)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorClickCounts(t *testing.T) {
	convey.Convey("ClickCounts", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ClickCounts(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
