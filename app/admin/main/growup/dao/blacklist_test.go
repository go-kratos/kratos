package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoListBlacklist(t *testing.T) {
	convey.Convey("ListBlacklist", t, func(ctx convey.C) {
		var (
			query = "id > 0"
			from  = int(0)
			limit = int(0)
			sort  = "-id"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(context.Background(), "INSERT INTO av_black_list(av_id, mid, ctype) VALUES(1993, 1001, 0)")
			list, total, err := d.ListBlacklist(query, from, limit, sort)
			ctx.Convey("Then err should be nil.list,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvIncomeStatis(t *testing.T) {
	convey.Convey("GetAvIncomeStatis", t, func(ctx convey.C) {
		var (
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(context.Background(), "INSERT INTO av_income_statis(av_id, total_income) VALUE(1000, 100)")
			avIncome, err := d.GetAvIncomeStatis(query)
			ctx.Convey("Then err should be nil.avIncome should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avIncome, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBlacklist(t *testing.T) {
	convey.Convey("UpdateBlacklist", t, func(ctx convey.C) {
		var (
			avID   = int64(1993)
			ctype  = int(0)
			update = make(map[string]interface{})
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(context.Background(), "DELETE FROM av_black_list WHERE av_id = 1993")
			Exec(context.Background(), "INSERT INTO av_black_list(av_id, mid, ctype) VALUES(1993, 1001, 0) ON DUPLICATE KEY UPDATE ctype = 0")
			update["ctype"] = 1
			err := d.UpdateBlacklist(avID, ctype, update)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
