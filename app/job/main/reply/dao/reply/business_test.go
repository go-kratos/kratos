package reply

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyListBusiness(t *testing.T) {
	convey.Convey("ListBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			business, err := d.Business.ListBusiness(c)
			ctx.Convey("Then err should be nil.business should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(business, convey.ShouldNotBeNil)
			})
		})
	})
}
