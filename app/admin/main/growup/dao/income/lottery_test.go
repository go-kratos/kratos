package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetLotteryIncome(t *testing.T) {
	convey.Convey("GetLotteryIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = ""
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(2000)
			typ   = int(5)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO lottery_av_income(av_id, mid, income, date) VALUES(10010, 1001, 100, '2018-11-11') ONDUPLICATE KEY UPDATE date = '2018-11-11'")
			avs, err := d.GetLotteryIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}
