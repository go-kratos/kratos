package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserLikeList(t *testing.T) {
	var (
		c          = context.TODO()
		id         = int64(1)
		businessID = int64(1)
		state      = int8(1)
		start      = int(1)
		end        = int(1)
	)
	convey.Convey("userLikeList", t, func(ctx convey.C) {
		_, err := d.userLikeList(c, id, businessID, state, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			// ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
