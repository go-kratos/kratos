package history

import (
	"context"
	"testing"

	"go-common/app/interface/main/history/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestHistorykeyHistory(t *testing.T) {
	convey.Convey("keyHistory", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyHistory(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistorykeyIndex(t *testing.T) {
	convey.Convey("keyIndex", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyIndex(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistorykeySwitch(t *testing.T) {
	convey.Convey("keySwitch", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keySwitch(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryExpireIndexCache(t *testing.T) {
	convey.Convey("ExpireIndexCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.ExpireIndexCache(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryExpireCache(t *testing.T) {
	convey.Convey("ExpireCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.ExpireCache(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryexpireCache(t *testing.T) {
	convey.Convey("expireCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.expireCache(c, key)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryIndexCache(t *testing.T) {
	convey.Convey("IndexCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.IndexCache(c, mid, start, end)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryIndexCacheByTime(t *testing.T) {
	convey.Convey("IndexCacheByTime", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			start = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.IndexCacheByTime(c, mid, start)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistorySetShadowCache(t *testing.T) {
	convey.Convey("SetShadowCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			value = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetShadowCache(c, mid, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryShadowCache(t *testing.T) {
	convey.Convey("ShadowCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			value, err := d.ShadowCache(c, mid)
			ctx.Convey("Then err should be nil.value should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryCacheMap(t *testing.T) {
	convey.Convey("CacheMap", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			amap, err := d.CacheMap(c, mid)
			ctx.Convey("Then err should be nil.amap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(amap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryCache(t *testing.T) {
	convey.Convey("Cache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{14771787}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.Cache(c, mid, aids)
			ctx.Convey("Then err should be nil.amap,miss should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryClearCache(t *testing.T) {
	convey.Convey("ClearCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ClearCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryDelCache(t *testing.T) {
	convey.Convey("DelCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{14771787}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCache(c, mid, aids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryAddCache(t *testing.T) {
	convey.Convey("AddCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			h   = &model.History{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCache(c, mid, h)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryAddCacheMap(t *testing.T) {
	convey.Convey("AddCacheMap", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			hm  map[int64]*model.History
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMap(c, mid, hm)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryTrimCache(t *testing.T) {
	convey.Convey("TrimCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			limit = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TrimCache(c, mid, limit)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryPingRedis(t *testing.T) {
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
