package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeredisKey(t *testing.T) {
	convey.Convey("redisKey", t, func(ctx convey.C) {
		var (
			key = "1"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := redisKey(key)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRsSet(t *testing.T) {
	convey.Convey("RsSet", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = "111"
			value = "1"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.RsSet(c, key, value)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeRbSet(t *testing.T) {
	convey.Convey("RbSet", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = "111"
			value = []byte("1")
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.RbSet(c, key, value)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeRsGet(t *testing.T) {
	convey.Convey("RsGet", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "1"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
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
			key = "2"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.RsSetNX(c, key)
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
			key = "1"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.Incr(c, key)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeCreateSelection(t *testing.T) {
	convey.Convey("CreateSelection", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			aid   = int64(1)
			stage = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.CreateSelection(c, aid, stage)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeSelection(t *testing.T) {
	convey.Convey("Selection", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			aid   = int64(1)
			stage = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			yes, no, err := d.Selection(c, aid, stage)
			ctx.Convey("Then err should be nil.yes,no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
				ctx.So(yes, convey.ShouldNotBeNil)
			})
		})
	})
}
