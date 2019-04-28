package toview

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/history/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestToviewkey(t *testing.T) {
	convey.Convey("key", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := key(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewExpire(t *testing.T) {
	convey.Convey("Expire", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.Expire(c, mid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewCache(t *testing.T) {
	convey.Convey("Cache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			start = int(1)
			end   = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Cache(c, mid, start, end)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewCacheMap(t *testing.T) {
	convey.Convey("CacheMap", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMap(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewCntCache(t *testing.T) {
	convey.Convey("CntCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CntCache(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewClearCache(t *testing.T) {
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

func TestToviewDelCaches(t *testing.T) {
	convey.Convey("DelCaches", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCaches(c, mid, aids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewAddCache(t *testing.T) {
	convey.Convey("AddCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			aid = int64(14771787)
			now = time.Now().Unix()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCache(c, mid, aid, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewAddCacheList(t *testing.T) {
	convey.Convey("AddCacheList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			views = []*model.ToView{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheList(c, mid, views)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewaddCache(t *testing.T) {
	convey.Convey("addCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = ""
			views = []*model.ToView{{Aid: 1477, Unix: 11}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.addCache(c, key, views)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewPingRedis(t *testing.T) {
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
