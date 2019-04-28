package like

import (
	"context"
	l "go-common/app/interface/main/activity/model/like"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeAddUserAchievment(t *testing.T) {
	convey.Convey("AddUserAchievment", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			userAchi = &l.ActLikeUserAchievement{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ID, err := d.AddUserAchievment(c, userAchi)
			ctx.Convey("Then err should be nil.ID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeUserAchievement(t *testing.T) {
	convey.Convey("UserAchievement", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserAchievement(c, sid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRawActUserAchieve(t *testing.T) {
	convey.Convey("RawActUserAchieve", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawActUserAchieve(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeActUserAchieveChange(t *testing.T) {
	convey.Convey("ActUserAchieveChange", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			award = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upID, err := d.ActUserAchieveChange(c, id, award)
			ctx.Convey("Then err should be nil.upID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upID, convey.ShouldNotBeNil)
			})
		})
	})
}
