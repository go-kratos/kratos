package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaospamKey(t *testing.T) {
	var (
		mid = int64(1234567)
		now = time.Now()
	)
	convey.Convey("spamKey", t, func(ctx convey.C) {
		p1 := spamKey(mid, now)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoIncrSpamCache(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(1234567)
		action = int32(5)
	)
	convey.Convey("IncrSpamCache", t, func(ctx convey.C) {
		err := d.IncrSpamCache(c, mid, action)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSpamCache(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(1234567)
		action = int32(5)
	)
	convey.Convey("SpamCache", t, func(ctx convey.C) {
		count, err := d.SpamCache(c, mid, action)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
