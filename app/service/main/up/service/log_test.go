package service

import (
	"context"
	"testing"

	"go-common/app/service/main/up/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicesendUpSpecialLog(t *testing.T) {
	convey.Convey("sendUpSpecialLog", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			opInfo = &UpSpecialLogInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.sendUpSpecialLog(c, opInfo)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServicefillGroupInfo(t *testing.T) {
	convey.Convey("fillGroupInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			up = &model.UpSpecialWithName{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.fillGroupInfo(c, up)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
