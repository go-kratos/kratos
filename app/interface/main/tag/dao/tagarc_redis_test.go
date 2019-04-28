package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/archive/api"
	"go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoarcKey(t *testing.T) {
	var (
		tid = int64(0)
	)
	convey.Convey("arcKey", t, func(ctx convey.C) {
		p1 := arcKey(tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoregionArcKey(t *testing.T) {
	var (
		rid = int32(0)
		tid = int64(0)
	)
	convey.Convey("regionArcKey", t, func(ctx convey.C) {
		p1 := regionArcKey(rid, tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoregionOriArcKey(t *testing.T) {
	var (
		rid = int32(0)
		tid = int64(0)
	)
	convey.Convey("regionOriArcKey", t, func(ctx convey.C) {
		p1 := regionOriArcKey(rid, tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaozrange(t *testing.T) {
	var (
		c     = context.TODO()
		key   = ""
		start = int(0)
		end   = int(0)
	)
	convey.Convey("zrange", t, func(ctx convey.C) {
		aids, count, err := d.zrange(c, key, start, end)
		ctx.Convey("Then err should be nil.aids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(aids, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoNewArcsCache(t *testing.T) {
	var (
		c     = context.TODO()
		tid   = int64(0)
		start = int(0)
		end   = int(0)
	)
	convey.Convey("NewArcsCache", t, func(ctx convey.C) {
		aids, count, err := d.NewArcsCache(c, tid, start, end)
		ctx.Convey("Then err should be nil.aids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(aids, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRemoveNewArcsCache(t *testing.T) {
	var (
		c    = context.TODO()
		aid  = int64(0)
		tids = int64(0)
	)
	convey.Convey("RemoveNewArcsCache", t, func(ctx convey.C) {
		err := d.RemoveNewArcsCache(c, aid, tids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelNewArcsCache(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("DelNewArcsCache", t, func(ctx convey.C) {
		err := d.DelNewArcsCache(c, tid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRemTidArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		aid  = int64(10001)
		tids = []int64{1, 2, 3, 4}
	)
	convey.Convey("RemTidArcCache", t, func(ctx convey.C) {
		err := d.RemTidArcCache(c, aid, tids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDeleteNewArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		tid  = int64(0)
		aids = ""
	)
	convey.Convey("DeleteNewArcCache", t, func(ctx convey.C) {
		err := d.DeleteNewArcCache(c, tid, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddNewArcsCache(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
		as  = &api.Arc{}
	)
	convey.Convey("AddNewArcsCache", t, func(ctx convey.C) {
		err := d.AddNewArcsCache(c, tid, as)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddNewArcCache(t *testing.T) {
	var (
		c = context.TODO()
		a = &api.Arc{
			Aid:     10001,
			PubDate: time.Time(1544512542),
		}
		tids = int64(2)
	)
	convey.Convey("AddNewArcCache", t, func(ctx convey.C) {
		err := d.AddNewArcCache(c, a, tids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOriginRegionNewArcsCache(t *testing.T) {
	var (
		c     = context.TODO()
		rid   = int32(0)
		tid   = int64(0)
		start = int(0)
		end   = int(0)
	)
	convey.Convey("OriginRegionNewArcsCache", t, func(ctx convey.C) {
		aids, count, err := d.OriginRegionNewArcsCache(c, rid, tid, start, end)
		ctx.Convey("Then err should be nil.aids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(aids, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoExpireRegionNewArcsCache(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int32(0)
		tid = int64(0)
	)
	convey.Convey("ExpireRegionNewArcsCache", t, func(ctx convey.C) {
		ok, err := d.ExpireRegionNewArcsCache(c, rid, tid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoExpireOriginalNewestArcCache(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int32(0)
		tid = int64(0)
	)
	convey.Convey("ExpireOriginalNewestArcCache", t, func(ctx convey.C) {
		ok, err := d.ExpireOriginalNewestArcCache(c, rid, tid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRegionNewArcsCache(t *testing.T) {
	var (
		c     = context.TODO()
		rid   = int32(0)
		tid   = int64(0)
		start = int(0)
		end   = int(0)
	)
	convey.Convey("RegionNewArcsCache", t, func(ctx convey.C) {
		aids, count, err := d.RegionNewArcsCache(c, rid, tid, start, end)
		ctx.Convey("Then err should be nil.aids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(aids, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoAddRegionNewArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		rid  = int32(0)
		arc  = &api.Arc{}
		tids = int64(0)
	)
	convey.Convey("AddRegionNewArcCache", t, func(ctx convey.C) {
		err := d.AddRegionNewArcCache(c, rid, arc, tids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRemoveRegionNewArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		rid  = int32(0)
		arc  = &api.Arc{}
		tids = int64(0)
	)
	convey.Convey("RemoveRegionNewArcCache", t, func(ctx convey.C) {
		err := d.RemoveRegionNewArcCache(c, rid, arc, tids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRegionNewestArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		rid  = int32(33)
		tid  = int64(2)
		aids = []int64{10001}
	)
	convey.Convey("RegionNewestArcCache", t, func(ctx convey.C) {
		exist, none, err := d.RegionNewestArcCache(c, rid, tid, aids)
		ctx.Convey("Then err should be nil.exist,none should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(none, convey.ShouldNotBeNil)
			ctx.So(exist, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDeleteRegionNewArcsCache(t *testing.T) {
	var (
		c    = context.TODO()
		tid  = int64(2)
		rid  = int32(33)
		arcs = []*api.Arc{
			{
				Aid:     10001,
				PubDate: time.Time(1544512542),
			},
		}
	)
	convey.Convey("DeleteRegionNewArcsCache", t, func(ctx convey.C) {
		err := d.DeleteRegionNewArcsCache(c, tid, rid, arcs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddRegionNewestArcCache(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int32(0)
		tid = int64(0)
		as  = []*api.Arc{}
	)
	convey.Convey("AddRegionNewestArcCache", t, func(ctx convey.C) {
		err := d.AddRegionNewestArcCache(c, rid, tid, as)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOriginalNewestArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		rid  = int32(33)
		tid  = int64(2)
		aids = []int64{10001}
	)
	convey.Convey("OriginalNewestArcCache", t, func(ctx convey.C) {
		exist, none, err := d.OriginalNewestArcCache(c, rid, tid, aids)
		ctx.Convey("Then err should be nil.exist,none should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(none, convey.ShouldNotBeNil)
			ctx.So(exist, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddOriginalNewestArcCache(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int32(33)
		tid = int64(2)
		as  = []*api.Arc{
			{
				Aid:     10001,
				PubDate: time.Time(1544512542),
			},
		}
	)
	convey.Convey("AddOriginalNewestArcCache", t, func(ctx convey.C) {
		err := d.AddOriginalNewestArcCache(c, rid, tid, as)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
