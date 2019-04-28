package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		businessID = int64(1)
		state      = int8(1)
		start      = int(0)
		end        = int(10)
	)
	convey.Convey("UserLikeList", t, func(ctx convey.C) {
		_, err := d.UserLikeList(c, mid, businessID, state, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			// ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
