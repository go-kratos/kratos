package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/workflow/model/search"

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

func TestDaoChallCountCache(t *testing.T) {
	convey.Convey("ChallCountCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ChallCountCache(c, uid)
			ctx.Convey("Then err should be nil.challCount should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpChallCountCache(t *testing.T) {
	convey.Convey("UpChallCountCache", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			challCount = &search.ChallCount{}
			uid        = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpChallCountCache(c, challCount, uid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
