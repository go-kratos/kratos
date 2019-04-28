package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/vipinfo/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyInfo(t *testing.T) {
	convey.Convey("keyInfo", t, func(ctx convey.C) {
		var (
			mid = _testMid
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyInfo(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyFrozen(t *testing.T) {
	convey.Convey("keyFrozen", t, func(ctx convey.C) {
		var (
			mid = _testMid
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyFrozen(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheInfo(t *testing.T) {
	convey.Convey("CacheInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = _testMid
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.CacheInfo(c, mid)
			ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheInfo(t *testing.T) {
	convey.Convey("AddCacheInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = _testMid
			v   = &model.VipUserInfo{Mid: _testMid, VipType: 1, VipStatus: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheInfo(c, mid, v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheInfos(t *testing.T) {
	convey.Convey("CacheInfos", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = _testMids
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			for _, v := range mids {
				err := d.DelInfoCache(c, v)
				ctx.So(err, convey.ShouldBeNil)
			}
			res, err := d.CacheInfos(c, mids)
			fmt.Println("item", res)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheInfos(t *testing.T) {
	convey.Convey("AddCacheInfos", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			vs := make(map[int64]*model.VipUserInfo)
			vs[_testMid] = &model.VipUserInfo{Mid: _testMid, VipStatus: 1}
			err := d.AddCacheInfos(c, vs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheVipFrozen(t *testing.T) {
	convey.Convey("CacheVipFrozen", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = _testMid
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			val, err := d.CacheVipFrozen(c, mid)
			fmt.Println("TestDaoCacheVipFrozen", val)
			ctx.Convey("Then err should be nil.val should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(val, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCacheVipFrozens(t *testing.T) {
	convey.Convey("CacheVipFrozens", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			val, err := d.CacheVipFrozens(context.Background(), _testMids)
			fmt.Println("TestDaoCacheVipFrozens", val)
			ctx.Convey("Then err should be nil.val should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(val, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheFrozen(t *testing.T) {
	convey.Convey("AddCacheFrozen", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheFrozen(c, 1540883325, 1)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelInfoCache(t *testing.T) {
	convey.Convey("DelInfoCache", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelInfoCache(c, _testMid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
