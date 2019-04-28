package dao

import (
	"context"
	"go-common/app/service/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserLikesKey(t *testing.T) {
	var (
		businessID = int64(1)
		mid        = int64(1)
		state      = int8(1)
	)
	convey.Convey("userLikesKey", t, func(ctx convey.C) {
		p1 := userLikesKey(businessID, mid, state)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCacheUserLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		state      = int8(1)
		start      = int(1)
		end        = int(10)
	)
	convey.Convey("CacheUserLikeList", t, func(ctx convey.C) {
		_, err := d.CacheUserLikeList(c, mid, businessID, state, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			// ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheUserLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		miss       = []*model.ItemLikeRecord{}
		businessID = int64(1)
		state      = int8(1)
	)
	convey.Convey("AddCacheUserLikeList", t, func(ctx convey.C) {
		err := d.AddCacheUserLikeList(c, mid, miss, businessID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoExpireUserLikesCache(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		state      = int8(1)
	)
	convey.Convey("ExpireUserLikesCache", t, func(ctx convey.C) {
		ok, err := d.ExpireUserLikesCache(c, mid, businessID, state)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserLikeExists(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		messageIDs = []int64{}
		state      = int8(1)
	)
	convey.Convey("UserLikeExists", t, func(ctx convey.C) {
		res, err := d.UserLikeExists(c, mid, businessID, messageIDs, state)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAppendCacheUserLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		item       = &model.ItemLikeRecord{}
		businessID = int64(1)
		state      = int8(1)
	)
	convey.Convey("AppendCacheUserLikeList", t, func(ctx convey.C) {
		err := d.AppendCacheUserLikeList(c, mid, item, businessID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelUserLikeCache(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		messageID  = int64(1)
		state      = int8(1)
	)
	convey.Convey("DelUserLikeCache", t, func(ctx convey.C) {
		err := d.DelUserLikeCache(c, mid, businessID, messageID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUserLikesCountCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		mid        = int64(1)
	)
	convey.Convey("UserLikesCountCache", t, func(ctx convey.C) {
		res, err := d.UserLikesCountCache(c, businessID, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
