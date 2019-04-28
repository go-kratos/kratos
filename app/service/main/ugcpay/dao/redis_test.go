package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddCacheAssetRelationState(t *testing.T) {
	convey.Convey("AddCacheAssetRelationState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(2233)
			otype = "archive"
			mid   = int64(46333)
			state = "paid"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheAssetRelationState(c, oid, otype, mid, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoassetRelationKey(t *testing.T) {
	convey.Convey("assetRelationKey", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := assetRelationKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoassetRelationField(t *testing.T) {
	convey.Convey("assetRelationField", t, func(ctx convey.C) {
		var (
			oid   = int64(2233)
			otype = "archive"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := assetRelationField(oid, otype)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheAssetRelationState(t *testing.T) {
	convey.Convey("CacheAssetRelationState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(2233)
			otype = "archive"
			mid   = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			state, err := d.CacheAssetRelationState(c, oid, otype, mid)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldEqual, "paid")
			})
		})
	})
}

func TestDaoDelCacheAssetRelationState(t *testing.T) {
	convey.Convey("DelCacheAssetRelationState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(2233)
			otype = "archive"
			mid   = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheAssetRelationState(c, oid, otype, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheAssetRelation(t *testing.T) {
	convey.Convey("DelCacheAssetRelation", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheAssetRelation(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
