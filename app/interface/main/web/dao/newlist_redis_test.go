package dao

import (
	"context"
	"testing"
	xtime "time"

	"go-common/app/service/main/archive/api"
	"go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyNl(t *testing.T) {
	convey.Convey("keyNl", t, func(ctx convey.C) {
		var (
			rid = int32(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyNl(rid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyNlBak(t *testing.T) {
	convey.Convey("keyNlBak", t, func(ctx convey.C) {
		var (
			rid = int32(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyNlBak(rid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoNewListCache(t *testing.T) {
	convey.Convey("NewListCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int32(0)
			tp    = int8(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, count, err := d.NewListCache(c, rid, tp, start, end)
			ctx.Convey("Then err should be nil.arcs,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoNewListBakCache(t *testing.T) {
	convey.Convey("NewListBakCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int32(0)
			tp    = int8(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, count, err := d.NewListBakCache(c, rid, tp, start, end)
			ctx.Convey("Then err should be nil.arcs,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoSetNewListCache(t *testing.T) {
	convey.Convey("SetNewListCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int32(0)
			tp    = int8(0)
			arcs  = []*api.Arc{}
			count = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetNewListCache(c, rid, tp, arcs, count)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaofrom(t *testing.T) {
	convey.Convey("from", t, func(ctx convey.C) {
		var (
			i = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := from(i)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocombine(t *testing.T) {
	convey.Convey("combine", t, func(ctx convey.C) {
		var (
			pubdate = time.Time(xtime.Now().Unix())
			count   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := combine(pubdate, count)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
