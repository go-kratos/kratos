package reply

import (
	"context"
	"go-common/app/interface/main/reply/conf"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewCreditDao(t *testing.T) {
	convey.Convey("NewCreditDao", t, func(ctx convey.C) {
		var (
			c = conf.Conf
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := NewCreditDao(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyIsCreditUser(t *testing.T) {
	convey.Convey("IsCreditUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.CreditUser.IsCreditUser(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
