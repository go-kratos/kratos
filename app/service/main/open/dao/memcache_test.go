package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAppkeyCache(t *testing.T) {
	convey.Convey("AppkeyCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.AppkeyCache(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetAppkeyCache(t *testing.T) {
	convey.Convey("SetAppkeyCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			newData map[string]string
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetAppkeyCache(c, newData)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
