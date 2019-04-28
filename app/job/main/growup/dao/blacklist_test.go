package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"

	"go-common/app/job/main/growup/model"
)

func TestDaoListBlacklist(t *testing.T) {
	convey.Convey("ListBlacklist", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "av_id=1"
			from  = int(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_black_list(av_id, mid, reason, ctype, has_signed, nickname) VALUES (1,2,1,3,1,'test') ON DUPLICATE KEY UPDATE mid=VALUES(mid),reason=VALUES(reason),has_signed=VALUES(has_signed),nickname=VALUES(nickname)")
			backlists, err := d.ListBlacklist(c, query, from, limit)
			ctx.Convey("Then err should be nil.backlists should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(backlists, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetExecuteOrder(t *testing.T) {
	convey.Convey("GetExecuteOrder", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			startTime = time.Now()
			endTime   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			executeOrders, err := d.GetExecuteOrder(c, startTime, endTime)
			ctx.Convey("Then err should not be nil.executeOrders should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(executeOrders, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetLastCtime(t *testing.T) {
	convey.Convey("GetLastCtime", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			reason = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctime, err := d.GetLastCtime(c, reason)
			ctx.Convey("Then err should be nil.ctime should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ctime, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBlacklistBatch(t *testing.T) {
	convey.Convey("AddBlacklistBatch", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			blacklist = []*model.Blacklist{
				&model.Blacklist{
					ID:        int64(10),
					AvID:      int64(100),
					MID:       int64(3),
					Reason:    2,
					CType:     1,
					HasSigned: 1,
					Nickname:  "test",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.AddBlacklistBatch(c, blacklist)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetHasSignUpInfo(t *testing.T) {
	convey.Convey("GetHasSignUpInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int(0)
			limit  = int(100)
			m      = make(map[int64]string)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.GetHasSignUpInfo(c, offset, limit, m)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
