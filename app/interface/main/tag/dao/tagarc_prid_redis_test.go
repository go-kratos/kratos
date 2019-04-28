package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/archive/api"
	"go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotagPridKey(t *testing.T) {
	var (
		tid  = int64(1833)
		prid = int64(32)
	)
	convey.Convey("tagPridKey", t, func(ctx convey.C) {
		p1 := tagPridKey(tid, prid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoZrangeTagPridArc(t *testing.T) {
	var (
		c     = context.Background()
		tid   = int64(0)
		prid  = int64(0)
		start = int(0)
		end   = int(0)
	)
	convey.Convey("ZrangeTagPridArc", t, func(ctx convey.C) {
		aids, count, err := d.ZrangeTagPridArc(c, tid, prid, start, end)
		ctx.Convey("Then err should be nil.aids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(aids, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoAddTagPridArcCache(t *testing.T) {
	var (
		c    = context.Background()
		tids = []int64{1, 2, 3, 4}
		prid = int64(32)
		as   = &api.Arc{
			Aid:     10001,
			PubDate: time.Time(1544512542),
		}
	)
	convey.Convey("AddTagPridArcCache", t, func(ctx convey.C) {
		err := d.AddTagPridArcCache(c, tids, prid, as)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddTagPridCache(t *testing.T) {
	var (
		c    = context.Background()
		tids = []int64{1, 2, 3, 4}
		prid = int64(32)
		as   = &api.Arc{
			Aid:     10001,
			PubDate: time.Time(1544512542),
		}
	)
	convey.Convey("AddTagPridCache", t, func(ctx convey.C) {
		err := d.AddTagPridCache(c, tids, prid, as)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRemoveTagPridArcCache(t *testing.T) {
	var (
		c    = context.Background()
		tids = []int64{1, 2, 3, 4}
		prid = int64(32)
		aid  = int64(10001)
	)
	convey.Convey("RemoveTagPridArcCache", t, func(ctx convey.C) {
		err := d.RemoveTagPridArcCache(c, tids, prid, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoexpireTagArcCache(t *testing.T) {
	var (
		c    = context.Background()
		tid  = int64(0)
		prid = int64(0)
	)
	convey.Convey("expireTagArcCache", t, func(ctx convey.C) {
		ok, err := d.expireTagArcCache(c, tid, prid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}
