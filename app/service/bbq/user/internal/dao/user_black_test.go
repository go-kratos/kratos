package dao

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/bbq/user/internal/model"
	xtime "go-common/library/time"
	"testing"
	"time"
)

func TestDaoTxAddUserBlack(t *testing.T) {
	convey.Convey("TxAddUserBlack", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(88895104)
			upMid = int64(88895105)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxAddUserBlack(c, tx, mid, upMid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxCancelUserBlack(t *testing.T) {
	convey.Convey("TxCancelUserBlack", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(88895104)
			upMid = int64(88895105)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxCancelUserBlack(c, tx, mid, upMid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFetchBlackList(t *testing.T) {
	convey.Convey("FetchBlackList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upMid, err := d.FetchBlackList(c, mid)
			ctx.Convey("Then err should be nil.upMid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upMid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFetchPartBlackList(t *testing.T) {
	convey.Convey("FetchPartBlackList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(88895104)
			cursor model.CursorValue
			size   = int(10)
		)
		cursor.CursorID = model.MaxInt64
		cursor.CursorTime = xtime.Time(time.Now().Unix())
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			MID2IDMap, blackMIDs, err := d.FetchPartBlackList(c, mid, cursor, size)
			ctx.Convey("Then err should be nil.MID2IDMap,blackMIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(blackMIDs, convey.ShouldNotBeNil)
				ctx.So(MID2IDMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsBlack(t *testing.T) {
	convey.Convey("IsBlack", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			mid           = int64(88895104)
			candidateMIDs = []int64{88895105}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			MIDMap := d.IsBlack(c, mid, candidateMIDs)
			ctx.Convey("Then MIDMap should not be nil.", func(ctx convey.C) {
				ctx.So(MIDMap, convey.ShouldNotBeNil)
			})
		})
	})
}
