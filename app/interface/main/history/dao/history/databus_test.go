package history

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestHistoryPlayPro(t *testing.T) {
	convey.Convey("PlayPro", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "14771787"
			msg = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PlayPro(c, key, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryMerge(t *testing.T) {
	convey.Convey("Merge", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			msg = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Merge(c, mid, msg)
		})
	})
}

func TestHistoryexperiencePub(t *testing.T) {
	convey.Convey("experiencePub", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "14771787"
			msg = interface{}(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.experiencePub(c, key, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryProPub(t *testing.T) {
	convey.Convey("ProPub", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "14771787"
			msg = interface{}(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ProPub(c, key, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
