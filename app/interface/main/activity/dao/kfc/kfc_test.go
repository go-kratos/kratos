package kfc

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestKfcRawKfcCoupon(t *testing.T) {
	convey.Convey("RawKfcCoupon", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(4)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.RawKfcCoupon(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestKfcRawKfcCode(t *testing.T) {
	convey.Convey("RawKfcCode", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			code = "201812041203"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.RawKfcCode(c, code)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestKfcKfcCodeGiveOut(t *testing.T) {
	convey.Convey("KfcCodeGiveOut", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(3)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.KfcCodeGiveOut(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestKfcDeliver(t *testing.T) {
	convey.Convey("KfcCodeGiveOut", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			id  = int64(3)
			mid = int64(5874874)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.KfcDeliver(c, id, mid)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
