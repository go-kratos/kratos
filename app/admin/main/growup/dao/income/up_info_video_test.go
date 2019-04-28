package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetUpInfoNickname(t *testing.T) {
	convey.Convey("GetUpInfoNickname", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1993}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_info_video(mid, account_state) VALUES(1993, 3)")
			upInfo, err := d.GetUpInfoNickname(c, mids)
			ctx.Convey("Then err should be nil.upInfo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upInfo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetUpInfoNicknameByMID(t *testing.T) {
	convey.Convey("GetUpInfoNicknameByMID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(1993)
			table = "up_info_video"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nickname, err := d.GetUpInfoNicknameByMID(c, mid, table)
			ctx.Convey("Then err should be nil.nickname should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeTxUpdateUpInfoScore(t *testing.T) {
	convey.Convey("TxUpdateUpInfoScore", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			table = "up_info_video"
			score = int(5)
			mid   = int64(1993)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxUpdateUpInfoScore(tx, table, score, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
