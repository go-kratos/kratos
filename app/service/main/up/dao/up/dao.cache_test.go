package up

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpUp(t *testing.T) {
	convey.Convey("Up", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Up(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpSwitch(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("UpSwitch", t, func(ctx convey.C) {
		_, err := d.UpSwitch(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			//ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestUpUpInfoActive(t *testing.T) {
	convey.Convey("UpInfoActive", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UpInfoActive(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpsInfoActive(t *testing.T) {
	convey.Convey("UpsInfoActive", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UpsInfoActive(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
