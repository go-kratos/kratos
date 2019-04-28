package dao

import (
	"context"
	"go-common/app/service/main/history/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusinesses(t *testing.T) {
	convey.Convey("Businesses", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.Businesses(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeleteHistories(t *testing.T) {
	convey.Convey("DeleteHistories", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			bid       = int64(14771787)
			beginTime = time.Now()
			endTime   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.DeleteHistories(c, bid, beginTime, endTime)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddHistories(t *testing.T) {
	convey.Convey("AddHistories", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			hs = []*model.History{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddHistories(c, hs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDeleteUserHistories(t *testing.T) {
	convey.Convey("DeleteUserHistories", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			bid = int64(14771787)
			no  = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.DeleteUserHistories(c, mid, bid, no)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserHistories(t *testing.T) {
	convey.Convey("UserHistories", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(14771787)
			businessID = int64(3)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.UserHistories(c, mid, businessID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoEarlyHistory(t *testing.T) {
	convey.Convey("EarlyHistory", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(3)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.EarlyHistory(c, businessID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
