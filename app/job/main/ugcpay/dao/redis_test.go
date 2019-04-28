package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoassetRelationKey(t *testing.T) {
	convey.Convey("assetRelationKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
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
			oid   = int64(0)
			otype = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := assetRelationField(oid, otype)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheAssetRelationState(t *testing.T) {
	convey.Convey("DelCacheAssetRelationState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			otype = ""
			mid   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheAssetRelationState(c, oid, otype, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
