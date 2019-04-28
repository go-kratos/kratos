package data

import (
	"context"
	"go-common/app/interface/main/creative/model/data"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/memcache"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDatakeyBase(t *testing.T) {
	var (
		mid  = int64(1)
		date = ""
	)
	convey.Convey("keyBase", t, func(ctx convey.C) {
		p1 := keyBase(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyArea(t *testing.T) {
	var (
		mid  = int64(0)
		date = ""
	)
	convey.Convey("keyArea", t, func(ctx convey.C) {
		p1 := keyArea(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyTrend(t *testing.T) {
	var (
		mid  = int64(1)
		date = ""
	)
	convey.Convey("keyTrend", t, func(ctx convey.C) {
		p1 := keyTrend(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyRfd(t *testing.T) {
	var (
		mid  = int64(1)
		date = ""
	)
	convey.Convey("keyRfd", t, func(ctx convey.C) {
		p1 := keyRfd(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyRfm(t *testing.T) {
	var (
		mid  = int64(1)
		date = ""
	)
	convey.Convey("keyRfm", t, func(ctx convey.C) {
		p1 := keyRfm(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyAct(t *testing.T) {
	var (
		mid  = int64(0)
		date = ""
	)
	convey.Convey("keyAct", t, func(ctx convey.C) {
		p1 := keyAct(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyViewIncr(t *testing.T) {
	var (
		mid  = int64(1)
		ty   = ""
		date = ""
	)
	convey.Convey("keyViewIncr", t, func(ctx convey.C) {
		p1 := keyViewIncr(mid, ty, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyThirtyDayArchive(t *testing.T) {
	var (
		mid = int64(1)
		ty  = ""
	)
	convey.Convey("keyThirtyDayArchive", t, func(ctx convey.C) {
		p1 := keyThirtyDayArchive(mid, ty)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDatakeyThirtyDayArticle(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("keyThirtyDayArticle", t, func(ctx convey.C) {
		p1 := keyThirtyDayArticle(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDataViewerBaseCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		dt  = "dt"
		err error
	)
	convey.Convey("ViewerBaseCache1", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrExists)
		})
		defer connGuard.Unpatch()
		_, err = d.ViewerBaseCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("ViewerBaseCache2", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		_, err = d.ViewerBaseCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddViewerBaseCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		dt  = ""
		res map[string]*data.ViewerBase
	)
	convey.Convey("AddViewerBaseCache", t, func(ctx convey.C) {
		err := d.AddViewerBaseCache(c, mid, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataViewerAreaCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		dt  = "dt"
	)
	convey.Convey("ViewerAreaCache1", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrExists)
		})
		defer connGuard.Unpatch()
		_, err := d.ViewerAreaCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("ViewerAreaCache2", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		_, err := d.ViewerAreaCache(c, mid, dt)
		ctx.Convey("ShouldBeNil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddViewerAreaCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
		res map[string]map[string]int64
	)
	convey.Convey("AddViewerAreaCache", t, func(ctx convey.C) {
		err := d.AddViewerAreaCache(c, mid, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataTrendCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
	)
	convey.Convey("TrendCache", t, func(ctx convey.C) {
		_, err := d.TrendCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddTrendCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
		res map[string]*data.ViewerTrend
	)
	convey.Convey("AddTrendCache", t, func(ctx convey.C) {
		err := d.AddTrendCache(c, mid, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataRelationFansDayCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
	)
	convey.Convey("RelationFansDayCache", t, func(ctx convey.C) {
		_, err := d.RelationFansDayCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddRelationFansDayCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
		res map[string]map[string]int
	)
	convey.Convey("AddRelationFansDayCache", t, func(ctx convey.C) {
		err := d.AddRelationFansDayCache(c, mid, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataRelationFansMonthCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		dt  = ""
	)
	convey.Convey("RelationFansMonthCache", t, func(ctx convey.C) {
		_, err := d.RelationFansMonthCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddRelationFansMonthCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
		res map[string]map[string]int
	)
	convey.Convey("AddRelationFansMonthCache", t, func(ctx convey.C) {
		err := d.AddRelationFansMonthCache(c, mid, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataViewerActionHourCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
	)
	convey.Convey("ViewerActionHourCache", t, func(ctx convey.C) {
		_, err := d.ViewerActionHourCache(c, mid, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddViewerActionHourCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
		res map[string]*data.ViewerActionHour
	)
	convey.Convey("AddViewerActionHourCache", t, func(ctx convey.C) {
		err := d.AddViewerActionHourCache(c, mid, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataViewerIncrCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ty  = ""
		dt  = ""
	)
	convey.Convey("ViewerIncrCache", t, func(ctx convey.C) {
		_, err := d.ViewerIncrCache(c, mid, ty, dt)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddViewerIncrCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ty  = ""
		dt  = ""
		res = &data.ViewerIncr{}
	)
	convey.Convey("AddViewerIncrCache", t, func(ctx convey.C) {
		err := d.AddViewerIncrCache(c, mid, ty, dt, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataThirtyDayArchiveCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ty  = ""
	)
	convey.Convey("ThirtyDayArchiveCache", t, func(ctx convey.C) {
		_, err := d.ThirtyDayArchiveCache(c, mid, ty)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddThirtyDayArchiveCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ty  = ""
		res = []*data.ThirtyDay{}
	)
	convey.Convey("AddThirtyDayArchiveCache", t, func(ctx convey.C) {
		err := d.AddThirtyDayArchiveCache(c, mid, ty, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataThirtyDayArticleCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("ThirtyDayArticleCache", t, func(ctx convey.C) {
		_, err := d.ThirtyDayArticleCache(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataAddThirtyDayArticleCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		res = []*artmdl.ThirtyDayArticle{}
	)
	convey.Convey("AddThirtyDayArticleCache", t, func(ctx convey.C) {
		err := d.AddThirtyDayArticleCache(c, mid, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
