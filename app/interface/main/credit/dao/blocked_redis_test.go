package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoblockIndexKey(t *testing.T) {
	convey.Convey("blockIndexKey", t, func(convCtx convey.C) {
		var (
			otype = int8(0)
			btype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := blockIndexKey(otype, btype)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedIdxCache(t *testing.T) {
	convey.Convey("BlockedIdxCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			otype = int8(0)
			btype = int8(0)
			start = int(0)
			end   = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ids, err := d.BlockedIdxCache(c, otype, btype, start, end)
			convCtx.Convey("Then err should be nil.ids should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ids, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExpireBlockedIdx(t *testing.T) {
	convey.Convey("ExpireBlockedIdx", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			otype = int8(0)
			btype = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.ExpireBlockedIdx(c, otype, btype)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLoadBlockedIdx(t *testing.T) {
	convey.Convey("LoadBlockedIdx", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			otype = int8(0)
			btype = int8(0)
			infos = []*model.BlockedInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.LoadBlockedIdx(c, otype, btype, infos)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
