package dao

import (
	"context"
	"go-common/app/service/main/identify-game/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetAccessCache(t *testing.T) {
	var (
		c   = context.Background()
		key = "123456"
		res = &model.AccessInfo{Mid: 1, AppID: 1, Token: "ashiba", CreateAt: 123, UserID: "12", Name: "mdzz", Expires: 1577811661, Permission: "all"}
	)
	convey.Convey("SetAccessCache", t, func(ctx convey.C) {
		err := d.SetAccessCache(c, key, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})

	res.Expires = -1
	convey.Convey("SetAccessCacheExpires", t, func(ctx convey.C) {
		err := d.SetAccessCache(c, key, res)
		ctx.Convey("Then err should be error.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAccessCache(t *testing.T) {
	var (
		c   = context.Background()
		key = "123456"
	)
	convey.Convey("AccessCache", t, func(ctx convey.C) {
		res, err := d.AccessCache(c, key)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelAccessCache(t *testing.T) {
	var (
		c   = context.Background()
		key = "123456"
	)
	convey.Convey("DelAccessCache", t, func(ctx convey.C) {
		err := d.DelAccessCache(c, key)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
