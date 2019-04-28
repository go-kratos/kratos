package like

import (
	"context"
	l "go-common/app/interface/main/activity/model/like"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeRawLikeMissionBuff(t *testing.T) {
	convey.Convey("RawLikeMissionBuff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ID, err := d.RawLikeMissionBuff(c, sid, mid)
			ctx.Convey("Then err should be nil.ID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeMissionGroupAdd(t *testing.T) {
	convey.Convey("MissionGroupAdd", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			group = &l.MissionGroup{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			misID, err := d.MissionGroupAdd(c, group)
			ctx.Convey("Then err should be nil.misID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(misID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRawMissionGroupItems(t *testing.T) {
	convey.Convey("RawMissionGroupItems", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			lids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawMissionGroupItems(c, lids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}
