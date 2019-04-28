package dao

import (
	"context"
	"go-common/app/service/main/spy/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaouserInfoCacheKey(t *testing.T) {
	convey.Convey("userInfoCacheKey", t, func(ctx convey.C) {
		var (
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := userInfoCacheKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserInfoCache(t *testing.T) {
	convey.Convey("UserInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UserInfoCache(c, mid)
			ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddUserInfoCache(t *testing.T) {
	convey.Convey("AddUserInfoCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ui = &model.UserInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddUserInfoCache(c, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserInfoCache(t *testing.T) {
	convey.Convey("SetUserInfoCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ui = &model.UserInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserInfoCache(c, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelInfoCache(t *testing.T) {
	convey.Convey("DelInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelInfoCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
