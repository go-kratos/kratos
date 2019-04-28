package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserCard(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2233)
	)
	convey.Convey("UserCard", t, func(ctx convey.C) {
		res, err := d.UserCard(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserCards(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{2233}
	)
	convey.Convey("UserCards", t, func(ctx convey.C) {
		res, err := d.UserCards(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserProfile(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2233)
	)
	convey.Convey("UserProfile", t, func(ctx convey.C) {
		res, err := d.UserProfile(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
