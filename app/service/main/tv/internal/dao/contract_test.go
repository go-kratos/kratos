package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tv/internal/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserContractByMid(t *testing.T) {
	convey.Convey("UserContractByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			uc, err := d.UserContractByMid(c, mid)
			ctx.Convey("Then err should be nil.uc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(uc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserContractByContractId(t *testing.T) {
	convey.Convey("UserContractByContractId", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			contractId = "Wx45678934567893456789"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			uc, err := d.UserContractByContractId(c, contractId)
			ctx.Convey("Then err should be nil.uc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(uc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDeleteUserContract(t *testing.T) {
	convey.Convey("TxDeleteUserContract", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			id    = int32(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxDeleteUserContract(c, tx, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxInsertUserContract(t *testing.T) {
	convey.Convey("TxInsertUserContract", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			uc    = &model.UserContract{
				Mid:        27515308,
				ContractId: "Wx45678934567893456789",
				OrderNo:    "T234567890456789",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TxInsertUserContract(c, tx, uc)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
