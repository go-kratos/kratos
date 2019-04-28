package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLock(t *testing.T) {
	convey.Convey("Lock", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "LOCK:ORDER:12345"
			val = "123456789"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Lock(c, key, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUnlock(t *testing.T) {
	convey.Convey("Unlock", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "LOCK:ORDER:12345"
			val = "123456789"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Unlock(c, key, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
