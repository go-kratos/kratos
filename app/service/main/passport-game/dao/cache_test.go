package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTokenPBCache(t *testing.T) {
	var (
		c   = context.TODO()
		key = "123456"
	)
	convey.Convey("TokenPBCache", t, func(ctx convey.C) {
		res, err := d.TokenPBCache(c, key)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoInfoPBCache(t *testing.T) {
	var (
		c   = context.TODO()
		key = "123456"
	)
	convey.Convey("InfoPBCache", t, func(ctx convey.C) {
		res, err := d.InfoPBCache(c, key)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
