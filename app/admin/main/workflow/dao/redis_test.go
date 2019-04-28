package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaopingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRedisRPOPCids(t *testing.T) {
	convey.Convey("RedisRPOPCids", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			business = int8(0)
			round    = int64(0)
			num      = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cids, err := d.RedisRPOPCids(c, business, round, num)
			ctx.Convey("Then err should be nil.cids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsOnline(t *testing.T) {
	convey.Convey("IsOnline", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			assigneeAdminID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			online, err := d.IsOnline(c, assigneeAdminID)
			ctx.Convey("Then err should be nil.online should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(online, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddOnline(t *testing.T) {
	convey.Convey("AddOnline", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			assigneeAdminID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddOnline(c, assigneeAdminID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelOnline(t *testing.T) {
	convey.Convey("DelOnline", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			assigneeAdminID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelOnline(c, assigneeAdminID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoListOnline(t *testing.T) {
	convey.Convey("ListOnline", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ListOnline(c)
			ctx.Convey("Then err should be nil.ids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLogInOutTime(t *testing.T) {
	convey.Convey("LogInOutTime", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			uids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.LogInOutTime(c, uids)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaofieldOnlineList(t *testing.T) {
	convey.Convey("fieldOnlineList", t, func(ctx convey.C) {
		var (
			assigneeAdminID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.fieldOnlineList(assigneeAdminID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyChallCount(t *testing.T) {
	convey.Convey("keyChallCount", t, func(ctx convey.C) {
		var (
			assigneeAdminID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.keyChallCount(assigneeAdminID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
