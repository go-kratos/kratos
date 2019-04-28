package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserInfoKey(t *testing.T) {
	convey.Convey("userInfoKey", t, func(ctx convey.C) {
		var (
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := userInfoKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopayParamKey(t *testing.T) {
	convey.Convey("payParamKey", t, func(ctx convey.C) {
		var (
			token = "TOKEN:34567345678"
		)
		d.AddCachePayParam(context.TODO(), token, &model.PayParam{})

		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := payParamKey(token)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
