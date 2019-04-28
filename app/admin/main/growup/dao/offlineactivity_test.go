package dao

import (
	"context"
	"go-common/app/admin/main/growup/dao/shell"
	"go-common/app/admin/main/growup/model/offlineactivity"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaogenerateDelimiter(t *testing.T) {
	convey.Convey("generateDelimiter", t, func(ctx convey.C) {
		var (
			delimiter = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := generateDelimiter(delimiter)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaostrlistToInt64List(t *testing.T) {
	convey.Convey("strlistToInt64List", t, func(ctx convey.C) {
		var (
			list    = []string{"11", "22 ", "vv ", "33 "}
			trimSet = " "
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := strlistToInt64List(list, trimSet)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotrimString(t *testing.T) {
	convey.Convey("trimString", t, func(ctx convey.C) {
		var (
			strlist = []string{"11", "22 ", "vv ", "33 "}
			trimset = " "
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := trimString(strlist, trimset)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoParseMidsFromString(t *testing.T) {
	convey.Convey("ParseMidsFromString", t, func(ctx convey.C) {
		var (
			str = "123"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := ParseMidsFromString(str)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOfflineActivityAddActivity(t *testing.T) {
	convey.Convey("OfflineActivityAddActivity", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &offlineactivity.AddActivityArg{
				Title:     "test",
				BonusType: 0,
				Memo:      "memo test",
				Link:      "http://12345.com",
				BonusList: []*offlineactivity.BonusInfo{{TotalMoney: 100, Mids: "1,2,3,"}},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.OfflineActivityAddActivity(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoShellCallbackUpdate(t *testing.T) {
	convey.Convey("ShellCallbackUpdate", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			result = &shell.OrderCallbackJSON{
				ThirdOrderNo: "1001",
				Status:       "SUCCESS",
			}
			msgid = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO offline_activity_shell_order(id, result_id, order_id, order_status) VALUES(1000, '1001', '1001', 'SUCCESS')")
			orderInfo, err := d.ShellCallbackUpdate(c, result, msgid)
			ctx.Convey("Then err should be nil.orderInfo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(orderInfo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOfflineActivityGetUpBonusResult(t *testing.T) {
	convey.Convey("OfflineActivityGetUpBonusResult", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			needCount bool
			limit     = int(10)
			offset    = int(1)
			query     = "id>?"
			args      = "0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upResult, totalCount, err := d.OfflineActivityGetUpBonusResult(c, needCount, limit, offset, query, args)
			ctx.Convey("Then err should be nil.upResult,totalCount should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(totalCount, convey.ShouldNotBeNil)
				ctx.So(upResult, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOfflineActivityGetUpBonusResultSelect(t *testing.T) {
	convey.Convey("OfflineActivityGetUpBonusResultSelect", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			selectQuery = "activity_id"
			query       = "id>?"
			args        = "0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upResult, err := d.OfflineActivityGetUpBonusResultSelect(c, selectQuery, query, args)
			ctx.Convey("Then err should be nil.upResult should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upResult, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOfflineActivityGetUpBonusByActivityResult(t *testing.T) {
	convey.Convey("OfflineActivityGetUpBonusByActivityResult", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			limit  = int(10)
			offset = int(0)
			mid    = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upResult, totalCount, err := d.OfflineActivityGetUpBonusByActivityResult(c, limit, offset, mid)
			ctx.Convey("Then err should be nil.upResult,totalCount should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(totalCount, convey.ShouldNotBeNil)
				ctx.So(upResult, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOfflineActivityGetDB(t *testing.T) {
	convey.Convey("OfflineActivityGetDB", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.OfflineActivityGetDB()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateActivityState(t *testing.T) {
	convey.Convey("UpdateActivityState", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			activityID = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.UpdateActivityState(c, activityID)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}
