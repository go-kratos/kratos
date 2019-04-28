package like

import (
	"context"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeUnDoMatchUsers(t *testing.T) {
	convey.Convey("UnDoMatchUsers", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			matchObjID = int64(1)
			limit      = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			list, err := d.UnDoMatchUsers(c, matchObjID, limit)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", list)
			})
		})
	})
}

func TestLikeUpMatchUserResult(t *testing.T) {
	convey.Convey("UpMatchUserResult", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			matchObjID = int64(1)
			mids       = []int64{111}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpMatchUserResult(c, matchObjID, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeMatchObjInfo(t *testing.T) {
	convey.Convey("MatchObjInfo", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			matchObjID = int64(37)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.MatchObjInfo(c, matchObjID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}
