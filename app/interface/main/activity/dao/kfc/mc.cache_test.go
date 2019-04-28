package kfc

import (
	"context"
	"testing"

	"go-common/app/interface/main/activity/model/kfc"

	"github.com/smartystreets/goconvey/convey"
)

func TestKfcCacheKfcCoupon(t *testing.T) {
	convey.Convey("CacheKfcCoupon", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(3)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CacheKfcCoupon(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestKfcAddCacheKfcCoupon(t *testing.T) {
	convey.Convey("AddCacheKfcCoupon", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			id  = int64(3)
			val = &kfc.BnjKfcCoupon{ID: 3, CouponCode: "201812041201", Mid: 2089809}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheKfcCoupon(c, id, val)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestKfcDelCacheKfcCoupon(t *testing.T) {
	convey.Convey("DelCacheKfcCoupon", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(3)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCacheKfcCoupon(c, id)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestKfcCacheKfcCode(t *testing.T) {
	convey.Convey("CacheKfcCode", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = "201812041201"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CacheKfcCode(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestKfcAddCacheKfcCode(t *testing.T) {
	convey.Convey("AddCacheKfcCode", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			id  = "201812041201"
			val = int64(3)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheKfcCode(c, id, val)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
