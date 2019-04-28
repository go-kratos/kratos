package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/vip/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoVipInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("VipInfo", t, func(ctx convey.C) {
		r, err := d.VipInfo(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("r should not be nil", func(ctx convey.C) {
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateVipTypeAndStatus(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		vipStatus = int32(0)
		vipType   = int32(0)
	)
	convey.Convey("UpdateVipTypeAndStatus", t, func(ctx convey.C) {
		ret, err := d.UpdateVipTypeAndStatus(c, mid, vipStatus, vipType)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("ret should not be nil", func(ctx convey.C) {
			ctx.So(ret, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelChangeHistoryCount(t *testing.T) {
	var (
		c   = context.TODO()
		arg = &model.ArgChangeHistory{}
	)
	convey.Convey("SelChangeHistoryCount", t, func(ctx convey.C) {
		count, err := d.SelChangeHistoryCount(c, arg)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("count should not be nil", func(ctx convey.C) {
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelChangeHistory(t *testing.T) {
	var (
		c   = context.TODO()
		arg = &model.ArgChangeHistory{}
	)
	convey.Convey("SelChangeHistory", t, func(ctx convey.C) {
		_, err := d.SelChangeHistory(c, arg)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdateChannelID(t *testing.T) {
	var (
		c            = context.TODO()
		mid          = int64(0)
		payChannelID = int32(0)
		ver          = int64(0)
		oldVer       = int64(0)
	)
	convey.Convey("TxUpdateChannelID", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		err = d.TxUpdateChannelID(tx, mid, payChannelID, ver, oldVer)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxDupUserDiscount(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(0)
		discountID = int64(0)
		orderNo    = ""
		status     = int8(0)
	)
	convey.Convey("TxDupUserDiscount", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		err = d.TxDupUserDiscount(tx, mid, discountID, orderNo, status)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdatePayType(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		payType = int8(0)
		ver     = int64(0)
		oldVer  = int64(0)
	)
	convey.Convey("UpdatePayType", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		a, err := d.UpdatePayType(tx, mid, payType, ver, oldVer)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("a should not be nil", func(ctx convey.C) {
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSelVipChangeHistory(t *testing.T) {
	var (
		c          = context.TODO()
		relationID = ""
	)
	convey.Convey("SelVipChangeHistory", t, func(ctx convey.C) {
		r, err := d.SelVipChangeHistory(c, relationID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("r should not be nil", func(ctx convey.C) {
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertVipChangeHistory(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipChangeHistory{}
	)
	convey.Convey("InsertVipChangeHistory", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		id, err := d.InsertVipChangeHistory(tx, r)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxSelVipUserInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(20606508)
	)
	convey.Convey("TxSelVipUserInfo", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		r, err := d.TxSelVipUserInfo(tx, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("r should not be nil", func(ctx convey.C) {
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxAddIosVipUserInfo(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipInfoDB{
			Mid:                  time.Now().Unix(),
			VipType:              1,
			VipPayType:           1,
			VipStatus:            0,
			VipStartTime:         xtime.Time(time.Now().Unix()),
			VipOverdueTime:       xtime.Time(time.Now().Unix()),
			AnnualVipOverdueTime: xtime.Time(time.Now().Unix()),
			VipRecentTime:        xtime.Time(time.Now().Unix()),
		}
	)
	convey.Convey("TxAddIosVipUserInfo", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		eff, err := d.TxAddIosVipUserInfo(tx, r)
		err = tx.Commit()
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("eff should not be nil", func(ctx convey.C) {
			ctx.So(eff, convey.ShouldNotBeNil)
		})

		err = d.InsertVipUserInfo(tx, r)
		ctx.Convey("Error should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateIosUserInfo(t *testing.T) {
	var (
		c       = context.TODO()
		iosTime = xtime.Time(time.Now().Unix())
		mid     = int64(0)
	)

	convey.Convey("TxUpdateIosUserInfo", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		eff, err := d.TxUpdateIosUserInfo(tx, iosTime, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("eff should not be nil", func(ctx convey.C) {
			ctx.So(eff, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateIosRenewUserInfo(t *testing.T) {
	var (
		c            = context.TODO()
		paychannelID = int64(1)
		ver          = int64(1)
		oldVer       = int64(123)
		mid          = int64(20606508)
		payType      = int8(1)
	)
	convey.Convey("TxUpdateIosRenewUserInfo", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		err = d.TxUpdateIosRenewUserInfo(tx, paychannelID, ver, oldVer, mid, payType)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateVipUserInfo(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipInfoDB{
			Mid:                  206065080,
			VipType:              1,
			VipPayType:           1,
			VipStatus:            0,
			VipStartTime:         xtime.Time(time.Now().Unix()),
			VipOverdueTime:       xtime.Time(time.Now().Unix()),
			AnnualVipOverdueTime: xtime.Time(time.Now().Unix()),
			VipRecentTime:        xtime.Time(time.Now().Unix()),
		}
		ver = int64(0)
	)
	convey.Convey("UpdateVipUserInfo", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		a, err := d.UpdateVipUserInfo(tx, r, ver)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("a should not be nil", func(ctx convey.C) {
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}
