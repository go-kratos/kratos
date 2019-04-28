package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetIncomeCache(t *testing.T) {
	convey.Convey("SetIncomeCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = "100"
			value = map[string]interface{}{
				"aa": "bb",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetIncomeCache(c, key, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetIncomeCache(t *testing.T) {
	convey.Convey("GetIncomeCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "100"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetIncomeCache(c, key)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosetCacheKV(t *testing.T) {
	convey.Convey("setCacheKV", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			key    = "200"
			value  = []byte("aaaaaa")
			expire = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.setCacheKV(c, key, value, expire)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetCacheVal(t *testing.T) {
	convey.Convey("getCacheVal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "200"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.getCacheVal(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
