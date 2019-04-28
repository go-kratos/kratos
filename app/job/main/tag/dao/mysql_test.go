package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoreplace(t *testing.T) {
	convey.Convey("replace", t, func(ctx convey.C) {
		var (
			name = "this is tag "
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := replace(name)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagInfo(t *testing.T) {
	convey.Convey("TagInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TagInfo(c, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRids(t *testing.T) {
	convey.Convey("Rids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpMap, err := d.Rids(c)
			ctx.Convey("Then err should be nil.rpMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHotMap(t *testing.T) {
	convey.Convey("HotMap", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.HotMap(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.hit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResources(t *testing.T) {
	convey.Convey("Resources", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(9222)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Resources(c, oid, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagResources(t *testing.T) {
	convey.Convey("TagResources", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TagResources(c, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
