package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAppeals(t *testing.T) {
	convey.Convey("Appeals", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			appeals, err := d.Appeals(c, ids)
			ctx.Convey("Then err should be nil.appeals should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(appeals, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxSetWeight(t *testing.T) {
	convey.Convey("SetWeight", t, func(ctx convey.C) {
		var (
			newWeight map[int64]int64
			tx        = d.WriteORM.Begin()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxSetWeight(tx, newWeight)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			tx.Rollback()
		})
	})
}

func TestDaoSetAppealAssignState(t *testing.T) {
	convey.Convey("SetAppealAssignState", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			ids         = []int64{1}
			assignState = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetAppealAssignState(c, ids, assignState)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLastEvent(t *testing.T) {
	convey.Convey("LastEvent", t, func(ctx convey.C) {
		var apID = int64(1)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			e, err := d.LastEvent(apID)
			ctx.Convey("Then err should be nil.e should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(e, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetAppealTransferState(t *testing.T) {
	convey.Convey("SetAppealTransferState", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			ids           = []int64{}
			transferState = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetAppealTransferState(c, ids, transferState)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
