package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCompareInfo(t *testing.T) {
	var (
		c   = context.Background()
		rid = int64(1)
	)
	convey.Convey("CompareInfo", t, func(ctx convey.C) {
		err := d.CompareInfo(c, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSendMail(t *testing.T) {
	var (
		c      = context.Background()
		rid    = int64(1)
		status = int64(1)
	)
	convey.Convey("SendMail", t, func(ctx convey.C) {
		err := d.SendMail(c, rid, status)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
