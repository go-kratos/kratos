package dao

import (
	"context"
	"go-common/app/job/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserLikesKey(t *testing.T) {
	convey.Convey("userLikesKey", t, func(convCtx convey.C) {
		var (
			businessID = int64(33)
			mid        = int64(2233)
			state      = int8(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := userLikesKey(businessID, mid, state)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireUserLikesCache(t *testing.T) {
	convey.Convey("ExpireUserLikesCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(0)
			businessID = int64(0)
			state      = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.ExpireUserLikesCache(c, mid, businessID, state)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAppendCacheUserLikeList(t *testing.T) {
	convey.Convey("AppendCacheUserLikeList", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(0)
			item       = &model.ItemLikeRecord{}
			businessID = int64(0)
			state      = int8(0)
			limit      = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AppendCacheUserLikeList(c, mid, item, businessID, state, limit)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddUserLikesCache(t *testing.T) {
	convey.Convey("AddUserLikesCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(0)
			businessID = int64(0)
			items      = []*model.ItemLikeRecord{}
			typ        = int8(0)
			limit      = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddUserLikesCache(c, mid, businessID, items, typ, limit)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelUserLikeCache(t *testing.T) {
	convey.Convey("DelUserLikeCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(0)
			businessID = int64(0)
			messageID  = int64(0)
			state      = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelUserLikeCache(c, mid, businessID, messageID, state)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
