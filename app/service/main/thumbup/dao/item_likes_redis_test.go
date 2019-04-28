package dao

import (
	"context"
	"go-common/app/service/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoitemLikesKey(t *testing.T) {
	var (
		businessID = int64(1)
		messageID  = int64(1)
		state      = int8(1)
	)
	convey.Convey("itemLikesKey", t, func(ctx convey.C) {
		p1 := itemLikesKey(businessID, messageID, state)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCacheItemLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		businessID = int64(1)
		state      = int8(1)
		start      = int(1)
		end        = int(1)
	)
	convey.Convey("CacheItemLikeList", t, func(ctx convey.C) {
		_, err := d.CacheItemLikeList(c, messageID, businessID, state, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			// ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheItemLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		miss       = []*model.UserLikeRecord{}
		businessID = int64(1)
		state      = int8(1)
	)
	convey.Convey("AddCacheItemLikeList", t, func(ctx convey.C) {
		err := d.AddCacheItemLikeList(c, messageID, miss, businessID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoExpireItemLikesCache(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		businessID = int64(1)
		state      = int8(1)
	)
	convey.Convey("ExpireItemLikesCache", t, func(ctx convey.C) {
		ok, err := d.ExpireItemLikesCache(c, messageID, businessID, state)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoItemLikeExists(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		businessID = int64(1)
		mids       = []int64{}
		state      = int8(1)
	)
	convey.Convey("ItemLikeExists", t, func(ctx convey.C) {
		res, err := d.ItemLikeExists(c, messageID, businessID, mids, state)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAppendCacheItemLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		item       = &model.UserLikeRecord{}
		businessID = int64(1)
		state      = int8(1)
	)
	convey.Convey("AppendCacheItemLikeList", t, func(ctx convey.C) {
		err := d.AppendCacheItemLikeList(c, messageID, item, businessID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelItemLikeCache(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		businessID = int64(1)
		mid        = int64(1)
		state      = int8(1)
	)
	convey.Convey("DelItemLikeCache", t, func(ctx convey.C) {
		err := d.DelItemLikeCache(c, messageID, businessID, mid, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoItemLikesCountCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		messageID  = int64(1)
	)
	convey.Convey("ItemLikesCountCache", t, func(ctx convey.C) {
		res, err := d.ItemLikesCountCache(c, businessID, messageID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
