package dao

import (
	"context"
	"go-common/library/database/sql"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetAbleCode(t *testing.T) {
	convey.Convey("GetAbleCode", t, func(ctx convey.C) {
		var (
			tx          = &sql.Tx{}
			batchCodeID = int64(0)
		)
		tx, err := d.StartTx(context.Background())
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(tx, convey.ShouldNotBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetAbleCode(tx, batchCodeID)
			ctx.Convey("Then err should be nil.code should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateCodeRelationID(t *testing.T) {
	convey.Convey("UpdateCodeRelationID", t, func(ctx convey.C) {
		var (
			code       = ""
			relationID = ""
			bmid       = int64(0)
		)
		tx, err := d.StartTx(context.Background())
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(tx, convey.ShouldNotBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateCodeRelationID(tx, code, relationID, bmid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelBatchCodeByID(t *testing.T) {
	convey.Convey("SelBatchCodeByID", t, func(ctx convey.C) {
		var (
			batchCodeID = int64(0)
		)
		tx, err := d.StartTx(context.Background())
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(tx, convey.ShouldNotBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelBatchCodeByID(tx, batchCodeID)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
