package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoexpCoinKey(t *testing.T) {
	var (
		mid = int64(0)
		day = int64(0)
	)
	convey.Convey("expCoinKey", t, func(ctx convey.C) {
		p1 := expCoinKey(mid, day)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoexpAddedKey(t *testing.T) {
	var (
		tp  = ""
		mid = int64(0)
		day = int64(0)
	)
	convey.Convey("expAddedKey", t, func(ctx convey.C) {
		p1 := expAddedKey(tp, mid, day)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoStatCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		day = int64(0)
	)
	convey.Convey("StatCache", t, func(ctx convey.C) {
		st, err := d.StatCache(c, mid, day)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("st should not be nil", func(ctx convey.C) {
			ctx.So(st, convey.ShouldNotBeNil)
		})
	})
}
