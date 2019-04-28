package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserChangeHistoryByID(t *testing.T) {
	convey.Convey("UserChangeHistoryByID", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int32(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			uch, err := d.UserChangeHistoryByID(c, id)
			ctx.Convey("Then err should be nil.uch should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(uch, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserChangeHistorysByMid(t *testing.T) {
	convey.Convey("UserChangeHistorysByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			pn  = int32(1)
			ps  = int32(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.UserChangeHistorysByMid(c, mid, pn, ps)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserChangeHistorysByMidAndCtime(t *testing.T) {
	convey.Convey("UserChangeHistorysByMidAndCtime", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(27515308)
			from = xtime.Time(0)
			to   = xtime.Time(time.Now().Unix())
			pn   = int32(1)
			ps   = int32(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.UserChangeHistorysByMidAndCtime(c, mid, from, to, pn, ps)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertUserChangeHistory(t *testing.T) {
	convey.Convey("TxInsertUserChangeHistory", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			uch   = &model.UserChangeHistory{
				Mid:        27515308,
				ChangeType: model.UserChangeTypeCancelContract,
				ChangeTime: xtime.Time(time.Now().Unix()),
				OrderNo:    "T12345678345678",
				Days:       31,
				OperatorId: "",
				Remark:     "test",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TxInsertUserChangeHistory(c, tx, uch)
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
