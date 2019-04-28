package pendant

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPendantredPointFlagKey(t *testing.T) {
	convey.Convey("redPointFlagKey", t, func(ctx convey.C) {
		var (
			mid = int64(123)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := redPointFlagKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantRedPointCache(t *testing.T) {
	convey.Convey("RedPointCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(123)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			pid, err := d.RedPointCache(c, mid)
			ctx.Convey("Then err should be nil.pid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantSetRedPointCache(t *testing.T) {
	convey.Convey("SetRedPointCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(123)
			pid = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRedPointCache(c, mid, pid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPendantDelRedPointCache(t *testing.T) {
	convey.Convey("DelRedPointCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(123)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRedPointCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
