package datadao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatadaodescHelper(t *testing.T) {
	convey.Convey("descHelper", t, func(ctx convey.C) {
		var (
			d = &CacheBaseLoader{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			key := descHelper(d)
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaonewCacheBaseLoader(t *testing.T) {
	convey.Convey("newCacheBaseLoader", t, func(ctx convey.C) {
		var (
			signID = int64(0)
			date   = time.Now()
			val    = interface{}(0)
			desc   = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := newCacheBaseLoader(signID, date, val, desc)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoKeyCacheBaseLoader(t *testing.T) {
	var (
		c = CacheBaseLoader{}
	)
	convey.Convey("Key", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			key := c.Key()
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoValueCacheBaseLoader(t *testing.T) {
	var (
		c = CacheBaseLoader{Val: 1}
	)
	convey.Convey("Value", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			value := c.Value()
			ctx.Convey("Then value should not be nil.", func(ctx convey.C) {
				ctx.So(value, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoExpireCacheBaseLoader(t *testing.T) {
	var (
		c = CacheBaseLoader{Val: 1}
	)
	convey.Convey("Expire", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := c.Expire()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoDesc(t *testing.T) {
	var (
		c = CacheBaseLoader{Val: 1}
	)
	convey.Convey("Desc", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := c.Desc()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoNewCacheMcnDataWithTp(t *testing.T) {
	convey.Convey("NewCacheMcnDataWithTp", t, func(ctx convey.C) {
		var (
			signID   = int64(0)
			date     = time.Now()
			tp       = ""
			val      = interface{}(0)
			desc     = ""
			loadFunc LoadFuncWithTp
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := NewCacheMcnDataWithTp(signID, date, tp, val, desc, loadFunc)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoKeyCacheMcnDataWithTp(t *testing.T) {
	var (
		s = cacheMcnDataWithTp{CacheBaseLoader: CacheBaseLoader{Val: 1}}
	)
	convey.Convey("Key", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			key := s.Key()
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoLoadValuecacheMcnDataWithTp(t *testing.T) {
	convey.Convey("LoadValue", t, func(ctx convey.C) {
		var (
			c = context.Background()
			s = cacheMcnDataWithTp{
				CacheBaseLoader: CacheBaseLoader{Val: 1},
				LoadFunc: func(c context.Context, signID int64, date time.Time, tp string) (res interface{}, err error) {
					return d.GetIndexSource(c, signID, date, tp)
				},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			value, err := s.LoadValue(c)
			ctx.Convey("Then err should be nil.value should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoNewCacheMcnDataSignID(t *testing.T) {
	convey.Convey("NewCacheMcnDataSignID", t, func(ctx convey.C) {
		var (
			signID   = int64(0)
			date     = time.Now()
			val      = interface{}(0)
			desc     = ""
			loadFunc LoadFuncOnlySign
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := NewCacheMcnDataSignID(signID, date, val, desc, loadFunc)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoKeycacheMcnDataSignID(t *testing.T) {
	var (
		s = cacheMcnDataSignID{CacheBaseLoader: CacheBaseLoader{Val: 1}}
	)
	convey.Convey("Key", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			key := s.Key()
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoLoadValuecacheMcnDataSignID(t *testing.T) {
	convey.Convey("LoadValue", t, func(ctx convey.C) {
		var (
			c = context.Background()
			s = cacheMcnDataSignID{
				CacheBaseLoader: CacheBaseLoader{Val: 1},
				LoadFunc: func(c context.Context, signID int64, date time.Time) (res interface{}, err error) {
					return d.GetMcnFans(c, signID, date)
				},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			value, err := s.LoadValue(c)
			ctx.Convey("Then err should be nil.value should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldNotBeNil)
			})
		})
	})
}
