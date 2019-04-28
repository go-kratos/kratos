package dao

import (
	"context"
	"go-common/app/job/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubStatDatabus(t *testing.T) {
	convey.Convey("PubStatDatabus", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			business = "archive"
			mid      = int64(2233)
			s        = &model.Stats{}
			upMid    = int64(333)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.PubStatDatabus(c, business, mid, s, upMid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
