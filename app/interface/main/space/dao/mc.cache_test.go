package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddCacheNotice(t *testing.T) {
	convey.Convey("AddCacheNotice", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(2222)
			val = &model.Notice{Notice: "2222"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheNotice(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheNotice(t *testing.T) {
	convey.Convey("CacheNotice", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheNotice(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheNotice(t *testing.T) {
	convey.Convey("DelCacheNotice", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheNotice(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheTopArc(t *testing.T) {
	convey.Convey("AddCacheTopArc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(2222)
			val = &model.AidReason{Aid: 2222, Reason: "2222"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheTopArc(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheTopArc(t *testing.T) {
	convey.Convey("CacheTopArc", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheTopArc(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheMasterpiece(t *testing.T) {
	convey.Convey("AddCacheMasterpiece", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(2222)
			val = &model.AidReasons{List: []*model.AidReason{{Aid: 2222, Reason: "2222"}}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMasterpiece(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheMasterpiece(t *testing.T) {
	convey.Convey("CacheMasterpiece", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMasterpiece(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheTheme(t *testing.T) {
	convey.Convey("AddCacheTheme", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(2222)
			val = &model.ThemeDetails{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheTheme(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheTheme(t *testing.T) {
	convey.Convey("CacheTheme", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheTheme(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheTheme(t *testing.T) {
	convey.Convey("DelCacheTheme", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheTheme(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheTopDynamic(t *testing.T) {
	convey.Convey("AddCacheTopDynamic", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(2222)
			val = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheTopDynamic(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheTopDynamic(t *testing.T) {
	convey.Convey("CacheTopDynamic", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheTopDynamic(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
