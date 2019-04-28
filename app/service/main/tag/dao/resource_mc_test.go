package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyResTag(t *testing.T) {
	var (
		oid = int64(28843596)
		typ = int32(3)
	)
	convey.Convey("keyResTag", t, func(ctx convey.C) {
		p1 := keyResTag(oid, typ)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddResTagCache(t *testing.T) {
	var (
		c   = context.Background()
		oid = int64(28843596)
		typ = int32(3)
		rs  = []*model.Resource{
			{
				Oid:  28843596,
				Tid:  1833,
				Type: 3,
			},
		}
	)
	convey.Convey("AddResTagCache", t, func(ctx convey.C) {
		err := d.AddResTagCache(c, oid, typ, rs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddResTagMapCaches(t *testing.T) {
	var (
		c     = context.Background()
		tp    = int32(3)
		rsMap = map[int64][]*model.Resource{
			28843596: {
				{
					Oid:  28843596,
					Tid:  1833,
					Type: 3,
				},
			},
		}
	)
	convey.Convey("AddResTagMapCaches", t, func(ctx convey.C) {
		err := d.AddResTagMapCaches(c, tp, rsMap)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoResTagCache(t *testing.T) {
	var (
		c   = context.Background()
		oid = int64(28843596)
		tp  = int32(3)
	)
	convey.Convey("ResTagCache", t, func(ctx convey.C) {
		_, err := d.ResTagCache(c, oid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoResTagMapCache(t *testing.T) {
	var (
		c   = context.Background()
		oid = int64(28843596)
		tp  = int32(3)
	)
	convey.Convey("ResTagMapCache", t, func(ctx convey.C) {
		_, err := d.ResTagMapCache(c, oid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoResTagMapCaches(t *testing.T) {
	var (
		c    = context.Background()
		oids = []int64{28843596}
		tp   = int32(0)
	)
	convey.Convey("ResTagMapCaches", t, func(ctx convey.C) {
		res, missed, err := d.ResTagMapCaches(c, oids, tp)
		ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missed, convey.ShouldHaveLength, 1)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelResTagCache(t *testing.T) {
	var (
		c   = context.Background()
		oid = int64(28843596)
		tp  = int32(3)
	)
	convey.Convey("DelResTagCache", t, func(ctx convey.C) {
		err := d.DelResTagCache(c, oid, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddResourceMapCaches(t *testing.T) {
	var (
		c      = context.Background()
		rsMaps = map[int64][]*model.Resource{
			28843596: {
				{
					Oid:  28843596,
					Tid:  1833,
					Type: 3,
				},
			},
		}
	)
	convey.Convey("AddResourceMapCaches", t, func(ctx convey.C) {
		err := d.AddResourceMapCaches(c, rsMaps)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaokeyResource(t *testing.T) {
	convey.Convey("keyResource", t, func(ctx convey.C) {
		var (
			oid = int64(28843596)
			tid = int64(1833)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyResource(oid, tid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddResourceCache(t *testing.T) {
	convey.Convey("AddResourceCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.Resource{
				Oid:  28843596,
				Tid:  1833,
				Type: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddResourceCache(c, r)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddResourcesCache(t *testing.T) {
	convey.Convey("AddResourcesCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rs = []*model.Resource{
				{
					Oid:  28843596,
					Tid:  1833,
					Type: 3,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddResourcesCache(c, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddResourceMapCache(t *testing.T) {
	convey.Convey("AddResourceMapCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rsMap = map[int64]*model.Resource{
				28843596: {
					Oid:  28843596,
					Tid:  1833,
					Type: 3,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddResourceMapCache(c, rsMap)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoResourceMapCache(t *testing.T) {
	convey.Convey("ResourceMapCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(28843596)
			typ  = int32(3)
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, missed, err := d.ResourceMapCache(c, oid, typ, tids)
			ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResourceCache(t *testing.T) {
	convey.Convey("ResourceCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceCache(c, oid, typ, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
