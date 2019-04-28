package dao

import (
	"context"
	"math"
	"testing"

	"go-common/app/job/main/ugcpay/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawOldElecTradeInfo(t *testing.T) {
	convey.Convey("RawOldElecTradeInfo", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "2016081812152789113137"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.RawOldElecTradeInfo(c, orderID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpsertElecUPRank(t *testing.T) {
	convey.Convey("TXUpsertElecUPRank", t, func(ctx convey.C) {
		var (
			c              = context.Background()
			tx, _          = d.BeginTranRank(c)
			ver            = int64(0)
			upMID          = int64(233)
			payMID         = int64(322)
			deltaPayAmount = int64(0)
			hidden         bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TXUpsertElecUPRank(c, tx, ver, upMID, payMID, deltaPayAmount, hidden)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertElecAVRank(t *testing.T) {
	convey.Convey("TXUpsertElecAVRank", t, func(ctx convey.C) {
		var (
			c              = context.Background()
			tx, _          = d.BeginTranRank(c)
			ver            = int64(0)
			avID           = int64(233)
			upMID          = int64(233)
			payMID         = int64(322)
			deltaPayAmount = int64(0)
			hidden         bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TXUpsertElecAVRank(c, tx, ver, avID, upMID, payMID, deltaPayAmount, hidden)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUpsertElecMessage(t *testing.T) {
	convey.Convey("UpsertElecMessage", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.DBElecMessage{
				ID:      233333,
				AVID:    233,
				UPMID:   233,
				PayMID:  322,
				Message: "ut",
				Replied: true,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpsertElecMessage(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpsertElecReply(t *testing.T) {
	convey.Convey("UpsertElecReply", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.DBElecReply{
				ID:     233334,
				MSGID:  233333,
				Reply:  "ut_reply",
				Hidden: false,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpsertElecReply(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOldElecMessageList(t *testing.T) {
	convey.Convey("OldElecMessageList", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startID = int64(0)
			limit   = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, list, err := d.OldElecMessageList(c, startID, limit)
			ctx.Convey("Then err should be nil.maxID,list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOldElecOrderList(t *testing.T) {
	convey.Convey("OldElecOrderList", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startID = int64(0)
			limit   = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, list, err := d.OldElecOrderList(c, startID, limit)
			ctx.Convey("Then err should be nil.maxID,list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestDaoElecAddSetting(t *testing.T) {
	convey.Convey("ElecAddSetting", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(46333)
			orValue = int32(0x01)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ElecAddSetting(c, math.MaxInt32, mid, orValue)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoElecDeleteSetting(t *testing.T) {
	convey.Convey("ElecDeleteSetting", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(46333)
			andValue = int32(0x7ffffffd)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ElecDeleteSetting(c, math.MaxInt32, mid, andValue)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
