package dao

import (
	"context"
	"testing"

	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpMcnSignStateOP(t *testing.T) {
	convey.Convey("UpMcnSignStateOP", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			state  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMcnSignStateOP(c, signID, state)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpMcnUpStateOP(t *testing.T) {
	convey.Convey("UpMcnUpStateOP", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			signUpID = int64(0)
			state    = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMcnUpStateOP(c, signUpID, state)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpMcnSignPayExpOP(t *testing.T) {
	convey.Convey("UpMcnSignPayExpOP", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			signPayID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMcnSignPayExpOP(c, signPayID)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddMcnDataSummary(t *testing.T) {
	convey.Convey("AddMcnDataSummary", t, func(ctx convey.C) {
		var (
			c                   = context.Background()
			mcnMid              = int64(0)
			signID              = int64(1)
			upCount             = int64(0)
			fansCountAccumulate = int64(0)
			genDate             xtime.Time
		)
		var _, err = d.db.Exec(c, "delete from mcn_data_summary where sign_id=? and generate_date='1970-01-01'", signID)
		if err != nil {
			t.Logf("err=%v", err)
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMcnDataSummary(c, mcnMid, signID, upCount, fansCountAccumulate, genDate)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoMcnSigns(t *testing.T) {
	convey.Convey("McnSigns", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mss, err := d.McnSigns(c)
			ctx.Convey("Then err should be nil.mss should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mss, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMcnUps(t *testing.T) {
	convey.Convey("McnUps", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int64(0)
			limit  = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.McnUps(c, offset, limit)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(ups), convey.ShouldBeGreaterThanOrEqualTo,0)
			})
		})
	})
}

func TestDaoMcnSignPays(t *testing.T) {
	convey.Convey("McnSignPays", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			sps, err := d.McnSignPayWarns(c)
			ctx.Convey("Then err should be nil.sps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(sps), convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoMcnSignMids(t *testing.T) {
	convey.Convey("McnSignMids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			msid, sids, err := d.McnSignMids(c)
			ctx.Convey("Then err should be nil.msid,sids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sids, convey.ShouldNotBeNil)
				ctx.So(msid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMcnUPCount(t *testing.T) {
	convey.Convey("McnUPCount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			signIDs = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mmc, err := d.McnUPCount(c, signIDs)
			ctx.Convey("Then err should be nil.mmc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mmc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMcnUPMids(t *testing.T) {
	convey.Convey("McnUPMids", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			signIDs = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mup, err := d.McnUPMids(c, signIDs)
			ctx.Convey("Then err should be nil.mup should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mup, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCrmUpMidsSum(t *testing.T) {
	convey.Convey("CrmUpMidsSum", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			upMids = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CrmUpMidsSum(c, upMids)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
