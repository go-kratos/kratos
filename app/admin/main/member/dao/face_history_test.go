package dao

import (
	xtime "go-common/library/time"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDaoreverse(t *testing.T) {
	convey.Convey("reverse", t, func(ctx convey.C) {
		var (
			s = "123"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := reverse(s)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "321")
			})
		})
	})
}

func TestDaorpad(t *testing.T) {
	convey.Convey("rpad", t, func(ctx convey.C) {
		var (
			s = "123"
			c = "a"
			l = int(4)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := rpad(s, c, l)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "123a")
			})
		})
	})
}

func TestDaomidKey(t *testing.T) {
	convey.Convey("midKey", t, func(ctx convey.C) {
		var (
			mid  = int64(1)
			time xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := midKey(mid, time)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaooperatorKey(t *testing.T) {
	convey.Convey("operatorKey", t, func(ctx convey.C) {
		var (
			operator = "admin"
			time     xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := operatorKey(operator, time)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoscanTimes(t *testing.T) {
	convey.Convey("scanTimes", t, func(ctx convey.C) {
		var (
			tDuration xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := scanTimes(tDuration)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFaceHistoryByMid(t *testing.T) {
	convey.Convey("FaceHistoryByMid", t, func(ctx convey.C) {
		// var (
		// 	c   = context.Background()
		// 	arg = &model.ArgFaceHistory{}
		// )
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			// p1, err := d.FaceHistoryByMid(c, arg)
			// ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			// 	ctx.So(err, convey.ShouldBeNil)
			// 	ctx.So(p1, convey.ShouldNotBeNil)
			// })
		})
	})
}

func TestDaoFaceHistoryByOP(t *testing.T) {
	convey.Convey("FaceHistoryByOP", t, func(ctx convey.C) {
		// var (
		// 	c   = context.Background()
		// 	arg = &model.ArgFaceHistory{}
		// )
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			// p1, err := d.FaceHistoryByOP(c, arg)
			// ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			// 	ctx.So(err, convey.ShouldBeNil)
			// 	ctx.So(p1, convey.ShouldNotBeNil)
			// })
		})
	})
}

func TestDaotoOPFaceRecord(t *testing.T) {
	convey.Convey("toOPFaceRecord", t, func(ctx convey.C) {
		var (
			res = &hrpc.Result{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := toOPFaceRecord(res)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotoMidFaceRecord(t *testing.T) {
	convey.Convey("toMidFaceRecord", t, func(ctx convey.C) {
		var (
			res = &hrpc.Result{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := toMidFaceRecord(res)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
