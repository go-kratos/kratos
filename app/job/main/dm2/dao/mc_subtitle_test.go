package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosubtitleKey(t *testing.T) {
	convey.Convey("subtitleKey", t, func(ctx convey.C) {
		var (
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.subtitleKey(oid, subtitleID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubtitleVideoKey(t *testing.T) {
	convey.Convey("subtitleVideoKey", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.subtitleVideoKey(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosubtitleDraftKey(t *testing.T) {
	convey.Convey("subtitleDraftKey", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
			mid = int64(0)
			lan = uint8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.subtitleDraftKey(oid, tp, mid, lan)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelVideoSubtitleCache(t *testing.T) {
	convey.Convey("DelVideoSubtitleCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelVideoSubtitleCache(c, oid, tp)
		})
	})
}

func TestDaoDelSubtitleDraftCache(t *testing.T) {
	convey.Convey("DelSubtitleDraftCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
			mid = int64(0)
			lan = uint8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelSubtitleDraftCache(c, oid, tp, mid, lan)
		})
	})
}

func TestDaoDelSubtitleCache(t *testing.T) {
	convey.Convey("DelSubtitleCache", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelSubtitleCache(c, oid, subtitleID)
		})
	})
}
