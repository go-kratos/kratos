package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChallByIDs(t *testing.T) {
	convey.Convey("ChallByIDs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cids = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ChallByIDs(c, cids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpDispatchStateByIDs(t *testing.T) {
	convey.Convey("UpDispatchStateByIDs", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			cids          = []int64{}
			dispatchState = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpDispatchStateByIDs(c, cids, dispatchState)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpDispatchStateAdminIDByIds(t *testing.T) {
	convey.Convey("UpDispatchStateAdminIDByIds", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			cids          = []int64{}
			dispatchState = int64(0)
			assignAdminid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpDispatchStateAdminIDByIds(c, cids, dispatchState, assignAdminid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
