package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88895342)
	)
	convey.Convey("UserInfo", t, func(ctx convey.C) {
		res, err := d.UserInfo(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserInfos(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{88895342}
	)
	convey.Convey("UserInfos", t, func(ctx convey.C) {
		res, err := d.UserInfos(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
