package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDaoLoginLogs(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(123)
		limit = int(1)
	)
	convey.Convey("LoginLogs", t, func(ctx convey.C) {
		res, err := d.LoginLogs(c, mid, limit)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorowKeyLoginLog(t *testing.T) {
	var (
		mid = int64(123)
		ts  = int64(1530846373)
	)
	convey.Convey("rowKeyLoginLog", t, func(ctx convey.C) {
		res := rowKeyLoginLog(mid, ts)
		ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoreverse(t *testing.T) {
	var (
		b = []byte("123")
	)
	convey.Convey("reverse", t, func(ctx convey.C) {
		reverse(b)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestDaoscanLoginLog(t *testing.T) {
	var (
		cells = []*hrpc.Cell{
			{Family: []byte("f"), Qualifier: []byte("mid"), Value: []byte("123")},
			{Family: []byte("f"), Qualifier: []byte("ts"), Value: []byte("1530846373")},
			{Family: []byte("f"), Qualifier: []byte("ip"), Value: []byte("2886731194")},
			{Family: []byte("f"), Qualifier: []byte("t"), Value: []byte("0")},
			{Family: []byte("f"), Qualifier: []byte("s"), Value: []byte("s")},
		}
	)
	convey.Convey("scanLoginLog", t, func(ctx convey.C) {
		res, err := scanLoginLog(cells)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
