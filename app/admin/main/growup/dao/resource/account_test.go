package resource

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestResourceMidByNickname(t *testing.T) {
	convey.Convey("MidByNickname", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "name"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mid, err := MidByNickname(c, name)
			ctx.Convey("Then err should be nil.mid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceNamesByMIDs(t *testing.T) {
	convey.Convey("NamesByMIDs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := NamesByMIDs(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceNameByMID(t *testing.T) {
	convey.Convey("NameByMID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			nickname, err := NameByMID(c, mid)
			ctx.Convey("Then err should be nil.nickname should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}
