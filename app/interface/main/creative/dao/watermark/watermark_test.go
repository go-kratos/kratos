package watermark

import (
	"context"
	"go-common/app/interface/main/creative/model/watermark"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWatermarkAddWaterMark(t *testing.T) {
	var (
		c = context.TODO()
		w = &watermark.Watermark{}
	)
	convey.Convey("AddWaterMark", t, func(ctx convey.C) {
		id, err := d.AddWaterMark(c, w)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestWatermarkUpWaterMark(t *testing.T) {
	var (
		c = context.TODO()
		w = &watermark.Watermark{}
	)
	convey.Convey("UpWaterMark", t, func(ctx convey.C) {
		rows, err := d.UpWaterMark(c, w)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestWatermarkWaterMark(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("WaterMark", t, func(ctx convey.C) {
		w, err := d.WaterMark(c, mid)
		ctx.Convey("Then err should be nil.w should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(w, convey.ShouldNotBeNil)
		})
	})
}
