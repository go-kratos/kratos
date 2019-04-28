package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeredisKey(t *testing.T) {
	convey.Convey("redisKey", t, func(ctx convey.C) {
		var (
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := redisKey(key)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRsGet(t *testing.T) {
	convey.Convey("RsGet", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RsGet(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRsSetNX(t *testing.T) {
	convey.Convey("RsSetNX", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RsSetNX(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRb(t *testing.T) {
	convey.Convey("Rb", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Rb(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeIncr(t *testing.T) {
	convey.Convey("Incr", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Incr(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeIncrby(t *testing.T) {
	convey.Convey("Incrby", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Incrby(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
