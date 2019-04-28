package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotagPridKey(t *testing.T) {
	convey.Convey("tagPridKey", t, func(ctx convey.C) {
		var (
			tid  = int64(9222)
			prid = int64(17)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := tagPridKey(tid, prid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddTagPridArcCache(t *testing.T) {
	convey.Convey("AddTagPridArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			arc  = &model.Archive{Aid: 29661790, PubTime: "2018-10-23 15:59:45"}
			prid = int64(17)
			tids = []int64{9222}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddTagPridArcCache(c, arc, prid, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemTagPridArcCache(t *testing.T) {
	convey.Convey("RemTagPridArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aid  = int64(29661790)
			prid = int64(17)
			tids = []int64{9222}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemTagPridArcCache(c, aid, prid, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoexpireCache(t *testing.T) {
	convey.Convey("expireCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tid  = int64(9222)
			prid = int64(17)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.expireCache(c, tid, prid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}
