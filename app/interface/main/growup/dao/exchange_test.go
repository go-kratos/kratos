package dao

import (
	"context"
	"go-common/app/interface/main/growup/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetGoodsByProductID(t *testing.T) {
	convey.Convey("GetGoodsByProductID", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			productID = "aaaaa"
			goodsType = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO creative_goods(ex_resource_id,discount,goods_type,is_display) values('aaaaa', 100, 1, 2)")
			goods, err := d.GetGoodsByProductID(c, productID, goodsType)
			ctx.Convey("Then err should be nil.goods should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(goods, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetDisplayGoods(t *testing.T) {
	convey.Convey("GetDisplayGoods", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			isDisplay = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {

			Exec(c, "INSERT INTO creative_goods(ex_resource_id,discount,goods_type,is_display) values('aaaaa', 100, 1, 2)")
			goods, err := d.GetDisplayGoods(c, isDisplay)
			ctx.Convey("Then err should be nil.goods should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(goods, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertGoodsOrder(t *testing.T) {
	convey.Convey("TxInsertGoodsOrder", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			o     = &model.GoodsOrder{
				MID:        253550886,
				OrderNo:    "DHY-20181122172812-1617825777068670976-253550886",
				GoodsType:  1,
				GoodsID:    "incentive-1mon",
				GoodsName:  "月度大会员",
				GoodsPrice: 100,
				GoodsCost:  100,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "DELETE FROM creative_order WHERE order_no = 'DHY-20181122172812-1617825777068670976-253550886'")
			rows, err := d.TxInsertGoodsOrder(tx, o)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoListVipProducts(t *testing.T) {
	convey.Convey("ListVipProducts", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(212312)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.ListVipProducts(c, mid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExchangeBigVIP(t *testing.T) {
	convey.Convey("ExchangeBigVIP", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(253550886)
			resourceID = int64(16)
			orderNo    = time.Now().UnixNano()
			remark     = "creative"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ExchangeBigVIP(c, mid, resourceID, orderNo, remark)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
