package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/account-recovery/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNickNameLog(t *testing.T) {
	var (
		c           = context.Background()
		nickNameReq = &model.NickNameReq{Mid: 111001254}
	)
	convey.Convey("NickNameLog", t, func(ctx convey.C) {
		res, err := d.NickNameLog(c, nickNameReq)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserBindLog(t *testing.T) {
	var (
		c        = context.Background()
		emailReq = &model.UserBindLogReq{
			Action: "emailBindLog",
			Mid:    23,
			Size:   100,
		}
	)
	convey.Convey("UserBindLog", t, func(ctx convey.C) {
		res, err := d.UserBindLog(c, emailReq)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
