package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaonotifyKey(t *testing.T) {
	var (
		mid = int64(4780461)
	)
	convey.Convey("notifyKey", t, func(ctx convey.C) {
		p1 := notifyKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddExplog(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(0)
		exp    = int64(0)
		toExp  = int64(0)
		oper   = ""
		reason = ""
		ip     = ""
	)
	convey.Convey("AddExplog", t, func(ctx convey.C) {
		err := d.AddExplog(c, mid, exp, toExp, oper, reason, ip)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoNotifyPurgeCache(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(0)
		action = ""
	)
	convey.Convey("NotifyPurgeCache", t, func(ctx convey.C) {
		p1 := d.NotifyPurgeCache(c, mid, action)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeNil)
		})
	})
}
