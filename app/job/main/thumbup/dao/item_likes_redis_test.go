package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/thumbup/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoitemLikesKey(t *testing.T) {
	convey.Convey("itemLikesKey", t, func(convCtx convey.C) {
		var (
			businessID = int64(33)
			messageID  = int64(5566)
			state      = int8(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := itemLikesKey(businessID, messageID, state)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireItemLikesCache(t *testing.T) {
	convey.Convey("ExpireItemLikesCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			messageID  = int64(5566)
			businessID = int64(33)
			state      = int8(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.ExpireItemLikesCache(c, messageID, businessID, state)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddItemLikesCache(t *testing.T) {
	convey.Convey("AddItemLikesCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(33)
			messageID  = int64(5566)
			typ        = int8(1)
			limit      = int(100)
			items      = []*model.UserLikeRecord{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddItemLikesCache(c, businessID, messageID, typ, limit, items)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAppendCacheItemLikeList(t *testing.T) {
	convey.Convey("AppendCacheItemLikeList", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			messageID  = int64(5566)
			item       = &model.UserLikeRecord{}
			businessID = int64(33)
			state      = int8(1)
			limit      = int(100)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AppendCacheItemLikeList(c, messageID, item, businessID, state, limit)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelItemLikeCache(t *testing.T) {
	convey.Convey("DelItemLikeCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			messageID  = int64(5566)
			businessID = int64(33)
			mid        = int64(2233)
			state      = int8(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelItemLikeCache(c, messageID, businessID, mid, state)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
