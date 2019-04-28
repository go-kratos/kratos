package dao

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/bbq/user/internal/model"
	xtime "go-common/library/time"
	"testing"
	"time"
)

func TestDaoTxAddUserFollow(t *testing.T) {
	convey.Convey("TxAddUserFollow", t, func(ctx convey.C) {
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
			num, err := d.TxAddUserFollow(c, tx, mid, upMid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxCancelUserFollow(t *testing.T) {
	convey.Convey("TxCancelUserFollow", t, func(ctx convey.C) {
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
			num, err := d.TxCancelUserFollow(c, tx, mid, upMid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFetchFollowList(t *testing.T) {
	convey.Convey("FetchFollowList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upMid, err := d.FetchFollowList(c, mid)
			ctx.Convey("Then err should be nil.upMid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upMid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFetchPartFollowList(t *testing.T) {
	convey.Convey("FetchPartFollowList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(88895104)
			cursor model.CursorValue
			size   = int(10)
		)
		cursor.CursorID = model.MaxInt64
		cursor.CursorTime = xtime.Time(time.Now().Unix())
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			MID2IDMap, followedMIDs, err := d.FetchPartFollowList(c, mid, cursor, size)
			ctx.Convey("Then err should be nil.MID2IDMap,followedMIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(followedMIDs, convey.ShouldNotBeNil)
				ctx.So(MID2IDMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsFollow(t *testing.T) {
	convey.Convey("IsFollow", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			mid           = int64(88895104)
			candidateMIDs = []int64{88895105}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			MIDMap := d.IsFollow(c, mid, candidateMIDs)
			ctx.Convey("Then MIDMap should not be nil.", func(ctx convey.C) {
				ctx.So(MIDMap, convey.ShouldNotBeNil)
			})
		})
	})
}
