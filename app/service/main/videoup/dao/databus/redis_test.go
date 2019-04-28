package databus

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatabusLock(t *testing.T) {
	var (
		c   = context.TODO()
		key = ""
		ttl = int(0)
	)
	convey.Convey("Lock", t, func(ctx convey.C) {
		gotLock, err := d.Lock(c, key, ttl)
		ctx.Convey("Then err should be nil.gotLock should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(gotLock, convey.ShouldNotBeNil)
		})
	})
}

func TestDatabusUnLock(t *testing.T) {
	var (
		c   = context.TODO()
		key = ""
	)
	convey.Convey("UnLock", t, func(ctx convey.C) {
		err := d.UnLock(c, key)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
