package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyBindByMid(t *testing.T) {
	convey.Convey("keyBindByMid", t, func(convCtx convey.C) {
		var (
			mid   = int64(0)
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := keyBindByMid(mid, appID)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

// func TestDaoCacheBindInfoByMid(t *testing.T) {
// 	convey.Convey("CacheBindInfoByMid", t, func(convCtx convey.C) {
// 		var (
// 			c     = context.Background()
// 			mid   = int64(0)
// 			appID = int64(0)
// 		)
// 		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
// 			v, err := d.CacheBindInfoByMid(c, mid, appID)
// 			convCtx.Convey("Then err should be nil.v should not be nil.", func(convCtx convey.C) {
// 				convCtx.So(err, convey.ShouldBeNil)
// 				convCtx.So(v, convey.ShouldNotBeNil)
// 			})
// 		})
// 	})
// }

func TestDaoAddCacheBindInfoByMid(t *testing.T) {
	convey.Convey("AddCacheBindInfoByMid", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			v     = &model.OpenBindInfo{}
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheBindInfoByMid(c, mid, v, appID)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelBindInfoCache(t *testing.T) {
	convey.Convey("DelBindInfoCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			appID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelBindInfoCache(c, mid, appID)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaodelCacheUnion(t *testing.T) {
	convey.Convey("delCacheUnion", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			key = "t_key"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.delCache(c, key)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheOpenInfoByOpenID(t *testing.T) {
	convey.Convey("AddCacheOpenInfoByOpenID", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			openID = "xx101"
			appID  = int64(32)
		)
		v := &model.OpenInfo{OpenID: openID, Mid: 101, AppID: appID}
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheOpenInfoByOpenID(c, openID, v, appID)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheOpenInfoByOpenID(t *testing.T) {
	convey.Convey("CacheOpenInfoByOpenID", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			openID = "xx101"
			appID  = int64(32)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			v, err := d.CacheOpenInfoByOpenID(c, openID, appID)
			fmt.Println("v--", v)
			convCtx.Convey("Then err should be nil.v should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(v, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelOpenInfoCache(t *testing.T) {
	convey.Convey("DelOpenInfoCache", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			openID = "xx101"
			appID  = int64(32)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelOpenInfoCache(c, openID, appID)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
