package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohit(t *testing.T) {
	var (
		id = int64(0)
	)
	convey.Convey("hit", t, func(ctx convey.C) {
		p1 := hit(id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagUserHit(t *testing.T) {
	var (
		id = int64(0)
	)
	convey.Convey("tagUserHit", t, func(ctx convey.C) {
		p1 := tagUserHit(id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserRelation(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("UserRelation", t, func(ctx convey.C) {
		f, err := d.UserRelation(c, mid, fid)
		ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(f, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserSetAchieveFlag(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		flag = uint64(0)
	)
	convey.Convey("UserSetAchieveFlag", t, func(ctx convey.C) {
		p1, err := d.UserSetAchieveFlag(c, mid, flag)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
