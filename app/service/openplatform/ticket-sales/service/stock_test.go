package service

import (
	"fmt"
	"testing"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStockLock(t *testing.T) {
	Convey("StockLock", t, func() {
		items := make([]*model.SkuCnt, 0)
		items = append(items, &model.SkuCnt{
			SkuID: 1,
			Count: 1,
		})
		err := svr.StockLock(ctx, consts.OpTypeOrder, 1, items, true)
		So(err, ShouldBeIn, ecode.TicketStockLack, nil)
	})
}

func TestStocksLock(t *testing.T) {
	Convey("StocksLock", t, func() {
		items := make([]*model.OrderStockCnt, 0)
		items = append(items, &model.OrderStockCnt{
			OrderID: 1,
			Count:   1,
		})
		err := svr.StocksLock(ctx, 2, items, true)
		So(err, ShouldBeIn, ecode.TicketStockLack, nil)
	})
}

func TestStocksUnlock(t *testing.T) {
	Convey("StocksUnlock", t, func() {
		err := svr.StocksUnlock(ctx, 1, 1, []int64{1, 2})
		So(err, ShouldBeIn, nil, ecode.TicketStockLogNotFound)
	})
}

func TestStockLockedDecr(t *testing.T) {
	Convey("StockLockedDecr", t, func() {
		items := make([]*model.SkuCnt, 0)
		items = append(items, &model.SkuCnt{
			SkuID: 1,
			Count: 1,
		})
		err := svr.StockLockedDecr(ctx, 1, 1, items)
		So(err, ShouldBeNil)
	})
}

func TestStockIncr(t *testing.T) {
	Convey("StockIncr", t, func() {
		items := make([]*model.SkuCnt, 0)
		items = append(items, &model.SkuCnt{
			SkuID: 1,
			Count: 1,
		})
		err := svr.StockIncr(ctx, 1, 1, items)
		So(err, ShouldBeNil)
	})
}

func TestStockInit(t *testing.T) {
	Convey("StockInit", t, func() {
		batches := make([]*model.Batch, 0)
		batches = append(batches, &model.Batch{
			ScreenID:      1,
			TicketPriceID: 1,
			SkAlert:       0,
			TotalStock:    10,
		})
		skuIds, err := svr.StockInit(ctx, 1, 1, batches)
		fmt.Println(skuIds, err)
		So(skuIds, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
