package academy

import (
	"context"
	"testing"

	"go-common/app/interface/main/creative/model/academy"

	"github.com/smartystreets/goconvey/convey"
)

func TestAcademyAddFeedBack(t *testing.T) {
	convey.Convey("AddFeedBack", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			f   = &academy.FeedBack{}
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.AddFeedBack(c, f, mid)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
