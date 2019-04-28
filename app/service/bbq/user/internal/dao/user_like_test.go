package dao

import (
	"context"
	"go-common/app/service/bbq/user/internal/model"
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTxAddUserLike(t *testing.T) {
	convey.Convey("TxAddUserLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(88895104)
			svid = int64(133)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxAddUserLike(tx, mid, svid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxCancelUserLike(t *testing.T) {
	convey.Convey("TxCancelUserLike", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(88895104)
			svid = int64(133)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxCancelUserLike(tx, mid, svid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUserLikeList(t *testing.T) {
	convey.Convey("GetUserLikeList", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(88895104)
			cursorNext bool
			cursor     model.CursorValue
			size       = int(10)
		)
		cursorNext = true
		cursor.CursorID = model.MaxInt64
		cursor.CursorTime = xtime.Time(time.Now().Unix())
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			likeSvs, err := d.GetUserLikeList(c, mid, cursorNext, cursor, size)
			ctx.Convey("Then err should be nil.likeSvs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(likeSvs, convey.ShouldNotBeNil)
			})
		})
	})
}
