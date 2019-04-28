package service

import (
	"context"
	"go-common/app/admin/main/growup/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicetradeDao(t *testing.T) {
	convey.Convey("tradeDao", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := s.tradeDao()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceSyncGoods(t *testing.T) {
	convey.Convey("SyncGoods", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			gt = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			eff, err := s.SyncGoods(c, gt)
			ctx.Convey("Then err should be nil.eff should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(eff, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGoodsList(t *testing.T) {
	convey.Convey("GoodsList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			from  = int(0)
			limit = int(20)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			total, res, err := s.GoodsList(c, from, limit)
			ctx.Convey("Then err should be nil.total,res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateGoodsInfo(t *testing.T) {
	convey.Convey("UpdateGoodsInfo", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			discount = int(100)
			ID       = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := s.UpdateGoodsInfo(c, discount, ID)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateGoodsDisplay(t *testing.T) {
	convey.Convey("UpdateGoodsDisplay", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			status = model.DisplayOff
			IDs    = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			eff, err := s.UpdateGoodsDisplay(c, status, IDs)
			ctx.Convey("Then err should be nil.eff should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(eff, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceonlineGoods(t *testing.T) {
	convey.Convey("onlineGoods", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			IDs = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := s.onlineGoods(c, IDs)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceofflineGoods(t *testing.T) {
	convey.Convey("offlineGoods", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			IDs = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := s.offlineGoods(c, IDs)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceOrderStatistics(t *testing.T) {
	convey.Convey("OrderStatistics", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.OrderQueryArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := s.OrderStatistics(c, arg)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceorderAll(t *testing.T) {
	convey.Convey("orderAll", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			where = "id<10"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			orders, err := s.orderAll(c, where)
			ctx.Convey("Then err should be nil.orders should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(orders, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceorderStatistics(t *testing.T) {
	convey.Convey("orderStatistics", t, func(ctx convey.C) {
		var (
			orders   = []*model.OrderInfo{}
			start    = time.Now()
			end      = time.Now().AddDate(0, 0, -1)
			timeType = model.Monthly
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := orderStatistics(orders, start, end, timeType)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceOrderExport(t *testing.T) {
	convey.Convey("OrderExport", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			arg   = &model.OrderQueryArg{}
			from  = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.OrderExport(c, arg, from, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceOrderList(t *testing.T) {
	convey.Convey("OrderList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			arg   = &model.OrderQueryArg{}
			from  = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			total, list, err := s.OrderList(c, arg, from, limit)
			ctx.Convey("Then err should be nil.total,list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicepreprocess(t *testing.T) {
	convey.Convey("preprocess", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.OrderQueryArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := preprocess(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeTrue)
			})
		})
	})
}

func TestServiceorderQueryStr(t *testing.T) {
	convey.Convey("orderQueryStr", t, func(ctx convey.C) {
		var (
			arg = &model.OrderQueryArg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := orderQueryStr(arg)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
