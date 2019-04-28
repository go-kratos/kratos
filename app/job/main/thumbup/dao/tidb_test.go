package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusiness(t *testing.T) {
	convey.Convey("Business", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.Business(c)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserLikes(t *testing.T) {
	convey.Convey("UserLikes", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(2233)
			businessID = int64(33)
			typ        = int8(1)
			limit      = int(100)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.UserLikes(c, mid, businessID, typ, limit)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoItemLikes(t *testing.T) {
	convey.Convey("ItemLikes", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(33)
			originID   = int64(7788)
			messageID  = int64(5566)
			typ        = int8(1)
			limit      = int(100)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.ItemLikes(c, businessID, originID, messageID, typ, limit)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateLikeState(t *testing.T) {
	convey.Convey("UpdateLikeState", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(2233)
			businessID = int64(33)
			originID   = int64(7788)
			messageID  = int64(5566)
			state      = int8(1)
			likeTime   = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.UpdateLikeState(c, mid, businessID, originID, messageID, state, likeTime)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateCounts(t *testing.T) {
	convey.Convey("UpdateCounts", t, func(convCtx convey.C) {
		var (
			c             = context.Background()
			businessID    = int64(33)
			originID      = int64(7788)
			messageID     = int64(5566)
			likeCounts    = int64(100)
			dislikeCounts = int64(100)
			upMid         = int64(999)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateCounts(c, businessID, originID, messageID, likeCounts, dislikeCounts, upMid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoStat(t *testing.T) {
	convey.Convey("Stat", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(33)
			originID   = int64(7788)
			messageID  = int64(5566)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.Stat(c, businessID, originID, messageID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
