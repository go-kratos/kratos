package dao

import (
	"context"
	"go-common/app/service/main/vip/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetWelfareList(t *testing.T) {
	convey.Convey("GetWelfareList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.ArgWelfareList{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetWelfareList(c, req)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountWelfare(t *testing.T) {
	convey.Convey("CountWelfare", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.ArgWelfareList{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountWelfare(c, req)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetRecommendWelfare(t *testing.T) {
	convey.Convey("GetRecommendWelfare", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.ArgWelfareList{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetRecommendWelfare(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountRecommendWelfare(t *testing.T) {
	convey.Convey("CountRecommendWelfare", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.ArgWelfareList{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountRecommendWelfare(c, req)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetWelfareTypeList(t *testing.T) {
	convey.Convey("GetWelfareTypeList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetWelfareTypeList(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetWelfareInfo(t *testing.T) {
	convey.Convey("GetWelfareInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetWelfareInfo(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetWelfareBatch(t *testing.T) {
	convey.Convey("GetWelfareBatch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			wid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetWelfareBatch(c, wid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetReceivedCode(t *testing.T) {
	convey.Convey("GetReceivedCode", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			wid = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetReceivedCode(c, wid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateWelfareCodeUser(t *testing.T) {
	convey.Convey("UpdateWelfareCodeUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int(0)
			mid = int64(0)
		)
		tx, _ := d.StartTx(c)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.UpdateWelfareCodeUser(c, tx, id, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateWelfareBatch(t *testing.T) {
	convey.Convey("UpdateWelfareBatch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			bid = int(0)
		)
		tx, err := d.StartTx(c)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.UpdateWelfareBatch(c, tx, bid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetWelfareCodeUnReceived(t *testing.T) {
	convey.Convey("GetWelfareCodeUnReceived", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			wid  = int64(0)
			bids = []string{}
		)
		bids = append(bids, "1")
		bids = append(bids, "2")
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetWelfareCodeUnReceived(c, wid, bids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetMyWelfare(t *testing.T) {
	convey.Convey("GetMyWelfare", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.GetMyWelfare(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddReceiveRedirectWelfare(t *testing.T) {
	convey.Convey("AddReceiveRedirectWelfare", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			wid = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddReceiveRedirectWelfare(c, wid, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountReceiveRedirectWelfare(t *testing.T) {
	convey.Convey("CountReceiveRedirectWelfare", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			wid = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountReceiveRedirectWelfare(c, wid, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertReceiveRecord(t *testing.T) {
	convey.Convey("InsertReceiveRecord", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			wid       = int64(0)
			mid       = int64(0)
			monthYear = int64(0)
		)
		tx, err := d.StartTx(c)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.InsertReceiveRecord(c, tx, wid, mid, monthYear)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
