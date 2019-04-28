package up

import (
	"context"
	"testing"

	"go-common/app/admin/main/mcn/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpTxAddMcnSignEntry(t *testing.T) {
	convey.Convey("TxAddMcnSignEntry", t, func(ctx convey.C) {
		var (
			tx, _      = d.BeginTran(context.Background())
			mcnMid     = int64(0)
			beginDate  = ""
			endDate    = ""
			permission = uint32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			lastID, err := d.TxAddMcnSignEntry(tx, mcnMid, beginDate, endDate, permission)
			t.Logf("last id=%d", lastID)
			defer tx.Rollback()
			ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpTxAddMcnSignPay(t *testing.T) {
	convey.Convey("TxAddMcnSignPay", t, func(ctx convey.C) {
		var (
			tx, _    = d.BeginTran(context.Background())
			mcnMid   = int64(0)
			signID   = int64(0)
			payValue = int64(0)
			dueDate  = ""
			note     = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxAddMcnSignPay(tx, mcnMid, signID, payValue, dueDate, note)
			defer tx.Rollback()
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMcnSignOP(t *testing.T) {
	convey.Convey("UpMcnSignOP", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			signID       = int64(0)
			state        = int8(0)
			rejectTime   xtime.Time
			rejectReason = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMcnSignOP(c, signID, state, rejectTime, rejectReason)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMcnUpOP(t *testing.T) {
	convey.Convey("UpMcnUpOP", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			signUpID     = int64(0)
			state        = int8(0)
			rejectTime   xtime.Time
			rejectReason = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMcnUpOP(c, signUpID, state, rejectTime, rejectReason)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpTxMcnUpPermitOP(t *testing.T) {
	convey.Convey("TxMcnUpPermitOP", t, func(ctx convey.C) {
		var (
			tx, _      = d.BeginTran(context.Background())
			signID     = int64(0)
			mcnMid     = int64(0)
			upMid      = int64(0)
			permission = uint32(0)
			upAuthLink = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TxMcnUpPermitOP(tx, signID, mcnMid, upMid, permission, upAuthLink)
			defer tx.Rollback()
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpTxUpPermitApplyOP(t *testing.T) {
	convey.Convey("TxUpPermitApplyOP", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			arg   = &model.McnUpPermissionApply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TxUpPermitApplyOP(tx, arg)
			defer tx.Rollback()
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnSigns(t *testing.T) {
	convey.Convey("McnSigns", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNSignStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ms, err := d.McnSigns(c, arg)
			ctx.Convey("Then err should be nil.ms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnSignTotal(t *testing.T) {
	convey.Convey("McnSignTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNSignStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.McnSignTotal(c, arg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpTotal(t *testing.T) {
	convey.Convey("McnUpTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.McnUpTotal(c, arg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpPermits(t *testing.T) {
	convey.Convey("McnUpPermits", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPPermitStateReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ms, err := d.McnUpPermits(c, arg)
			ctx.Convey("Then err should be nil.ms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpPermit(t *testing.T) {
	convey.Convey("McnUpPermit", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			m, err := d.McnUpPermit(c, id)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnUpPermitTotal(t *testing.T) {
	convey.Convey("McnUpPermitTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPPermitStateReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.McnUpPermitTotal(c, arg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
