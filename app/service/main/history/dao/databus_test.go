package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/history/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddHistoryMessage(t *testing.T) {
	var (
		c   = context.Background()
		k   = int(1)
		msg = []*model.Merge{{
			Mid:  1,
			Bid:  4,
			Time: 10000,
		}}
	)
	convey.Convey("AddHistoryMessage", t, func(ctx convey.C) {
		err := d.AddHistoryMessage(c, k, msg)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
