package dao

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/bbq/user/internal/model"
	xtime "go-common/library/time"
	"testing"
	"time"
)

func TestDaoTxAddUserFan(t *testing.T) {
	convey.Convey("TxAddUserFan", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(88895104)
			fanMid = int64(88895105)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxAddUserFan(tx, mid, fanMid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxCancelUserFan(t *testing.T) {
	convey.Convey("TxCancelUserFan", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(88895104)
			fanMid = int64(88895105)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxCancelUserFan(tx, mid, fanMid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsFan(t *testing.T) {
	convey.Convey("IsFan", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			mid           = int64(88895104)
			candidateMIDs = []int64{88895105}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			MIDMap := d.IsFan(c, mid, candidateMIDs)
			ctx.Convey("Then MIDMap should not be nil.", func(ctx convey.C) {
				ctx.So(MIDMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFetchPartFanList(t *testing.T) {
	convey.Convey("FetchPartFanList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(88895104)
			cursor model.CursorValue
			size   = int(10)
		)
		cursor.CursorID = model.MaxInt64
		cursor.CursorTime = xtime.Time(time.Now().Unix())
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			MID2IDMap, followedMIDs, err := d.FetchPartFanList(c, mid, cursor, size)
			ctx.Convey("Then err should be nil.MID2IDMap,followedMIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(followedMIDs, convey.ShouldNotBeNil)
				ctx.So(MID2IDMap, convey.ShouldNotBeNil)
			})
		})
	})
}
