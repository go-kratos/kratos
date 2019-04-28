package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddGoods(t *testing.T) {
	convey.Convey("AddGoods", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			fields = "ex_product_id, ex_resource_id, goods_type, is_display, discount"
			values = fmt.Sprintf("('%s', 1, 1, 1, 100)", time.Now().Format("2006-01-02 15:04:05"))
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.AddGoods(c, fields, values)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateGoods(t *testing.T) {
	convey.Convey("UpdateGoods", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			set   = "is_display=1"
			where = "is_deleted=0"
			IDs   = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := d.UpdateGoods(c, set, where, IDs)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGoodsList(t *testing.T) {
	convey.Convey("GoodsList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			where = "is_deleted=0"
			from  = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GoodsList(c, where, from, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGoodsCount(t *testing.T) {
	convey.Convey("GoodsCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			where = "is_deleted=0"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			total, err := d.GoodsCount(c, where)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOrderCount(t *testing.T) {
	convey.Convey("OrderCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			where = "is_deleted=0"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			total, err := d.OrderCount(c, where)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOrderList(t *testing.T) {
	convey.Convey("OrderList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			where = "is_deleted=0"
			from  = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.OrderList(c, where, from, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
