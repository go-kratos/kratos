package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/thumbup/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusinesses(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Businesses", t, func(ctx convey.C) {
		res, err := d.Businesses(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoLikeState(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		originID   = int64(1)
		messageID  = int64(1)
	)
	convey.Convey("LikeState", t, func(ctx convey.C) {
		res, err := d.LikeState(c, mid, businessID, originID, messageID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserLikeCount(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		mid        = int64(1)
		typ        = int8(1)
	)
	convey.Convey("UserLikeCount", t, func(ctx convey.C) {
		res, err := d.UserLikeCount(c, businessID, mid, typ)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawItemLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		messageID  = int64(1)
		businessID = int64(1)
		originID   = int64(1)
		state      = int8(1)
		start      = int(1)
		end        = int(1)
	)
	convey.Convey("RawItemLikeList", t, func(ctx convey.C) {
		_, err := d.RawItemLikeList(c, messageID, businessID, originID, state, start, end)
		ctx.Convey("Then err should be nil.res,all should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRawUserLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		state      = int8(1)
		start      = int(1)
		end        = int(1)
	)
	convey.Convey("RawUserLikeList", t, func(ctx convey.C) {
		_, err := d.RawUserLikeList(c, mid, businessID, state, start, end)
		ctx.Convey("Then err should be nil.res,all should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			// ctx.So(all, convey.ShouldNotBeNil)
			// ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoMessageStats(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		ids        = []int64{1}
	)
	convey.Convey("MessageStats", t, func(ctx convey.C) {
		res, err := d.MessageStats(c, businessID, ids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeEmpty)
		})
	})
}

func TestDaoOriginStats(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
	)
	convey.Convey("OriginStats", t, func(ctx convey.C) {
		res, err := d.OriginStats(c, businessID, originID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoStat(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
		messageID  = int64(1)
	)
	convey.Convey("Stat", t, func(ctx convey.C) {
		res, err := d.Stat(c, businessID, originID, messageID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotidbRawStats(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
		messageID  = int64(1)
	)
	convey.Convey("RawStats", t, func(ctx convey.C) {
		res, err := d.RawStats(c, businessID, originID, messageID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateCount(t *testing.T) {
	var (
		c             = context.TODO()
		businessID    = int64(1)
		originID      = int64(1)
		messageID     = int64(1)
		likeChange    = int64(1)
		dislikeChange = int64(1)
	)
	convey.Convey("UpdateCount", t, func(ctx convey.C) {
		err := d.UpdateCount(c, businessID, originID, messageID, likeChange, dislikeChange)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

// INSERT INTO `bilibili_likes`.`counts`(`id`, `mtime`, `ctime`, `business_id`, `origin_id`, `message_id`, `likes_count`, `dislikes_count`, `likes_change`, `dislikes_change`, `up_mid`) VALUES (0, '2018-11-14 17:25:40', '2018-11-03 12:06:27', 3, 0, 10099865, 100, 0, 0, 0, 8167601);
func Test_updateCounts(t *testing.T) {
	convey.Convey("get data", t, func() {
		c := context.Background()
		messageID := int64(10099865)
		bid := int64(3)
		oid := int64(0)
		stat, err := d.Stat(c, bid, oid, messageID)
		convey.So(err, convey.ShouldBeNil)
		args := [][2]int64{
			{1, 1},
			{-1000, -1000},
			{1, 0},
			{0, 1},
			{-1, 0},
			{-1, -1},
			{0, -1},
			{0, 0},
			{100, -10000},
			{100, 20},
			{-10, 20},
		}
		var l, dd int64
		for _, x := range args {
			l = x[0]
			dd = x[1]
			convey.Convey(fmt.Sprintf("like %v dislike %v", l, dd), func() {
				err := d.UpdateCounts(c, bid, oid, messageID, l, dd, 0)
				convey.So(err, convey.ShouldBeNil)
				nstat, err := d.Stat(c, bid, oid, messageID)
				convey.So(err, convey.ShouldBeNil)
				likes := stat.Likes + l
				if likes < 0 {
					likes = 0
				}
				dislikes := stat.Dislikes + dd
				if dislikes < 0 {
					dislikes = 0
				}
				convey.So(nstat.Likes, convey.ShouldEqual, likes)
				convey.So(nstat.Dislikes, convey.ShouldEqual, dislikes)
			})
		}
	})
}

func TestDaoUpdateUpMids(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(2)
	)
	data := []*model.UpMidsReq{
		{MessageID: 100, UpMid: 100},
		{MessageID: 200, UpMid: 200},
		{MessageID: 300, UpMid: 300},
	}
	convey.Convey("UpdateUpMids", t, func(ctx convey.C) {
		_, err := d.UpdateUpMids(c, businessID, data)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoItemHasLike(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
		messageID  = int64(1)
		mids       = []int64{1, 2, 3}
	)
	convey.Convey("ItemHasLike", t, func(ctx convey.C) {
		res, err := d.ItemHasLike(c, businessID, originID, messageID, mids, model.StateLike)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}
