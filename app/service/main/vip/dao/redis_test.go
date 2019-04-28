package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoaccLoginKey(t *testing.T) {
	convey.Convey("accLoginKey", t, func(ctx convey.C) {
		var (
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := accLoginKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldContainSubstring, fmt.Sprintf(_accLogin, mid))
			})
		})
	})
}

func TestDaoRemQueue(t *testing.T) {
	convey.Convey("RemQueue", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.RemQueue(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelRedisCache(t *testing.T) {
	convey.Convey("DelRedisCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelRedisCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLoginCount(t *testing.T) {
	convey.Convey("LoginCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.LoginCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoFrozenTime(t *testing.T) {
	convey.Convey("FrozenTime", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			frozenTime, err := d.FrozenTime(c, mid)
			ctx.Convey("Then err should be nil.frozenTime should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(frozenTime, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}
