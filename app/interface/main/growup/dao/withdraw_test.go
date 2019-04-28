package dao

import (
	"context"
	"go-common/app/interface/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpAccountCount(t *testing.T) {
	convey.Convey("GetUpAccountCount", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			dateVersion = "2017-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, has_sign_contract, total_unwithdraw_income, withdraw_date_version, is_deleted) VALUES(1001, 1, 100, '2018-10', 0) ON DUPLICATE KEY UPDATE has_sign_contract = 1, is_deleted = 0")
			count, err := d.GetUpAccountCount(c, dateVersion)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUpAccountByDate(t *testing.T) {
	convey.Convey("QueryUpAccountByDate", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			dateVersion = "2017-01-01"
			from        = int(0)
			limit       = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, has_sign_contract, total_unwithdraw_income, withdraw_date_version, is_deleted) VALUES(1001, 1, 100, '2018-10', 0) ON DUPLICATE KEY UPDATE has_sign_contract = 1, is_deleted = 0")
			upAccounts, err := d.QueryUpAccountByDate(c, dateVersion, from, limit)
			ctx.Convey("Then err should be nil.upAccounts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upAccounts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUpWithdrawByMID(t *testing.T) {
	convey.Convey("QueryUpWithdrawByMID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income_withdraw(mid, date_version) VALUES(1001, '2018-08')")
			upWithdraws, err := d.QueryUpWithdrawByMID(c, mid)
			ctx.Convey("Then err should be nil.upWithdraws should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upWithdraws, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUpWithdrawByMids(t *testing.T) {
	convey.Convey("QueryUpWithdrawByMids", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			mids        = []int64{1001}
			dateVersion = "2018-08"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income_withdraw(mid, date_version, is_deleted) VALUES(1001, '2018-08', 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			upWithdraws, err := d.QueryUpWithdrawByMids(c, mids, dateVersion)
			ctx.Convey("Then err should be nil.upWithdraws should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upWithdraws, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUpWithdrawRecord(t *testing.T) {
	convey.Convey("InsertUpWithdrawRecord", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			upWithdraw = &model.UpIncomeWithdraw{
				MID:         1002,
				DateVersion: "2018-09",
				State:       1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_income_withdraw WHERE mid = 1002")
			result, err := d.InsertUpWithdrawRecord(c, upWithdraw)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUpWithdrawByID(t *testing.T) {
	convey.Convey("QueryUpWithdrawByID", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(100001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income_withdraw(id, mid, date_version, is_deleted) VALUES(100001, 1006, '2018-08', 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			upWithdraw, err := d.QueryUpWithdrawByID(c, id)
			ctx.Convey("Then err should be nil.upWithdraw should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upWithdraw, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateUpWithdrawState(t *testing.T) {
	convey.Convey("TxUpdateUpWithdrawState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			id    = int64(100001)
			state = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income_withdraw(id, mid, date_version, is_deleted, state) VALUES(100001, 1001, '2018-08', 0, 1) ON DUPLICATE KEY UPDATE is_deleted = 0, state = 1")
			result, err := d.TxUpdateUpWithdrawState(tx, id, state)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateUpAccountWithdraw(t *testing.T) {
	convey.Convey("TxUpdateUpAccountWithdraw", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tx, _     = d.BeginTran(c)
			mid       = int64(1001)
			thirdCoin = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, has_sign_contract, total_unwithdraw_income, withdraw_date_version, is_deleted) VALUES(1001, 1, 100, '2018-10', 0) ON DUPLICATE KEY UPDATE has_sign_contract = 1, is_deleted = 0")
			result, err := d.TxUpdateUpAccountWithdraw(tx, mid, thirdCoin)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxQueryMaxUpWithdrawDateVersion(t *testing.T) {
	convey.Convey("TxQueryMaxUpWithdrawDateVersion", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			mid   = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income_withdraw(id, mid, date_version, is_deleted, state) VALUES(10001, 1001, '2018-08', 0, 1) ON DUPLICATE KEY UPDATE is_deleted = 0, state = 1")
			dateVersion, err := d.TxQueryMaxUpWithdrawDateVersion(tx, mid)
			ctx.Convey("Then err should be nil.dateVersion should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dateVersion, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxQueryUpAccountVersion(t *testing.T) {
	convey.Convey("TxQueryUpAccountVersion", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			mid   = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, has_sign_contract, total_unwithdraw_income, withdraw_date_version, is_deleted) VALUES(1001, 1, 100, '2018-10', 0) ON DUPLICATE KEY UPDATE has_sign_contract = 1, is_deleted = 0")
			version, err := d.TxQueryUpAccountVersion(tx, mid)
			ctx.Convey("Then err should be nil.version should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(version, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTxUpdateUpAccountUnwithdrawIncome(t *testing.T) {
	convey.Convey("TxUpdateUpAccountUnwithdrawIncome", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			tx, _       = d.BeginTran(c)
			mid         = int64(1005)
			dateVersion = "2018-10"
			version     = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, version, is_deleted) VALUES(1001, 1, 0) ON DUPLICATE KEY UPDATE version = 1, is_deleted = 0")
			result, err := d.TxUpdateUpAccountUnwithdrawIncome(tx, mid, dateVersion, version)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
