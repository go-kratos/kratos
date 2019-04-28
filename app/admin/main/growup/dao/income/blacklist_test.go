package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeListAvBlackList(t *testing.T) {
	convey.Convey("ListAvBlackList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			avID  = []int64{19930812}
			ctype = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_black_list(av_id, mid) VALUES(19930812,19930812)")
			avb, err := d.ListAvBlackList(c, avID, ctype)
			ctx.Convey("Then err should be nil.avb should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avb, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetAvBlackListByMID(t *testing.T) {
	convey.Convey("GetAvBlackListByMID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(19930812)
			typ = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_black_list(av_id, mid) VALUES(1002,19930812)")
			avb, err := d.GetAvBlackListByMID(c, mid, typ)
			ctx.Convey("Then err should be nil.avb should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avb, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeTxInsertAvBlackList(t *testing.T) {
	convey.Convey("TxInsertAvBlackList", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			val   = "(1001,1000,0,'test','szy',1,0)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxInsertAvBlackList(tx, val)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
