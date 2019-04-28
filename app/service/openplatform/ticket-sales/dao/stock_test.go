package dao

import (
	"strings"
	"testing"

	"go-common/app/service/openplatform/ticket-sales/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTxAddStockLog(t *testing.T) {
	Convey("TestTxAddStockLog", t, func() {
		tx, _ := d.BeginTx(ctx)
		err := d.TxAddStockLog(tx, &model.SKUStockLog{
			SKUID:  1,
			OpType: 1,
			SrcID:  1,
			Stock:  1,
		})
		tx.Rollback()
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			So(err, ShouldNotBeNil)
		} else {
			So(err, ShouldBeNil)
		}
	})
}

func TestSkuByItemID(t *testing.T) {
	Convey("TestSkuByItemID", t, func() {
		res, err := d.SkuByItemID(ctx, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestStocks(t *testing.T) {
	Convey("TestStocks", t, func() {
		res, err := d.Stocks(ctx, []int64{1}, false)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetSKUs(t *testing.T) {
	Convey("TestGetSKUs", t, func() {
		res, err := d.GetSKUs(ctx, []int64{1}, false)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
