package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyBuyFlag(t *testing.T) {
	convey.Convey("keyBuyFlag", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyBuyFlag(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyApplyFlag(t *testing.T) {
	convey.Convey("keyApplyFlag", t, func(ctx convey.C) {
		var (
			code = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyApplyFlag(code)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetBuyFlagCache(t *testing.T) {
	convey.Convey("SetBuyFlagCache", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			f   = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.SetBuyFlagCache(c, mid, f)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelBuyFlagCache(t *testing.T) {
	convey.Convey("DelBuyFlagCache", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelBuyFlagCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetApplyFlagCache(t *testing.T) {
	convey.Convey("SetApplyFlagCache", t, func(ctx convey.C) {
		var (
			code = ""
			f    = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.SetApplyFlagCache(c, code, f)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelApplyFlagCache(t *testing.T) {
	convey.Convey("DelApplyFlagCache", t, func(ctx convey.C) {
		var (
			code = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelApplyFlagCache(c, code)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
