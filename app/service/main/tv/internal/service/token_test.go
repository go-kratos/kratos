package service

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceTokenInfos(t *testing.T) {
	convey.Convey("TokenInfos", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tokens = []string{"TOKEN:34567345678"}
		)
		s.dao.AddCachePayParam(context.TODO(), tokens[0], &model.PayParam{})
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			tis, err := s.TokenInfos(c, tokens)
			ctx.Convey("Then err should be nil.tis should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tis, convey.ShouldNotBeNil)
			})
		})
	})
}
