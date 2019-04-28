package dao

import (
	"testing"
	"time"

	"context"
	"go-common/app/service/openplatform/ticket-item/model"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

// TestGetPriceIDs
func TestDao_GetPriceIDs(t *testing.T) {
	Convey("GetPriceIDs", t, func() {
		once.Do(startService)
		res, err := d.GetPriceIDs(1, 2)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

// TestCanDelTicket
func TestDao_CanDelTicket(t *testing.T) {
	Convey("CanDelTicket", t, func() {
		once.Do(startService)
		res := d.CanDelTicket(33, true)
		So(res, ShouldBeTrue)
	})
}

/**
// TestDelTicket
func TestDao_DelTicket(t *testing.T) {
	Convey("DelTicket", t, func() {
		once.Do(startService)
		tx := d.db.Begin()
		err := d.DelTicket(context.TODO(), tx, []int64{33}, []int64{25}, 2, true)
		tx.Commit()
		So(err, ShouldBeNil)
	})
}**/

// TestCreateOrUpdateTkPrice
func TestDao_CreateOrUpdateTkPrice(t *testing.T) {
	Convey("CreateOrUpdateTkPrice", t, func() {
		once.Do(startService)
		tx := d.db.Begin()
		priceInfo, err := d.CreateOrUpdateTkPrice(context.TODO(), tx, model.TicketPrice{
			ID:            0,
			ProjectID:     33,
			Desc:          "DESCCCC",
			Type:          1, // 单场票
			SaleType:      1,
			Color:         "#FF0000",
			BuyLimit:      0,
			PaymentMethod: 1,
			PaymentValue:  0,
			DescDetail:    "DESCDETAIL",
			IsSale:        1,   // 可售
			IsRefund:      -10, // 不可退
			OriginPrice:   -1,  // 未設置
			MarketPrice:   -1,
			SaleStart:     -62135596800, // 0000-00-00 00:00:00
			SaleEnd:       -62135596800,
		}, 0)
		spew.Dump(priceInfo)
		tx.Commit()
		So(priceInfo, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

/**
// TestInsertOrUpdateTkPass
func TestDao_InsertOrUpdateTkPass(t *testing.T) {
	Convey("InsertOrUpdateTkPass", t, func() {
		once.Do(startService)
		tx := d.db.Begin()
		priceInfo, err := d.InsertOrUpdateTkPass(context.TODO(), tx, model.TicketPrice{
			ID:            0,
			ProjectID:     33,
			Desc:          "DESCCCC",
			Type:          1, // 单场票
			SaleType:      1,
			Color:         "#FF0000",
			BuyLimit:      0,
			PaymentMethod: 1,
			PaymentValue:  0,
			DescDetail:    "DESCDETAIL",
			IsSale:        1,   // 可售
			IsRefund:      -10, // 不可退
			OriginPrice:   -1,  // 未設置
			MarketPrice:   -1,
			SaleStart:     -62135596800, // 0000-00-00 00:00:00
			SaleEnd:       -62135596800,
		}, 0)
		spew.Dump(priceInfo)
		tx.Commit()
		So(priceInfo, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}**/

// TestTicketPrice
func TestTicketPrice(t *testing.T) {
	var (
		tp  *model.TicketPrice
		tps []*model.TicketPrice
		err error
	)
	Convey("TxGetParentTicketPrice", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)
		tp, err = d.TxGetParentTicketPrice(c, tx, 13424)
		So(tp.Desc, ShouldEqual, "海尔测试票")
		So(err, ShouldBeNil)
		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)
	})

	var cnt int
	Convey("TxGetPriceSymbols", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)
		tps, err = d.TxGetPriceSymbols(c, tx, 1633)
		cnt = len(tps)
		So(err, ShouldBeNil)
		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)
	})

	var genID int64
	once.Do(startService)
	c := context.TODO()
	tx, _ := d.BeginTran(c)
	Convey("TxBatchAddTicketPrice", t, func() {
		genID, err = model.GetTicketIDFromBase()
		So(genID, ShouldNotEqual, 0)
		So(err, ShouldBeNil)
		tp.ParentID = tp.ID
		tp.ID = genID
		tp.ScreenID = 1633
		tp.Symbol = string(int(time.Now().Unix())%26 + 'A')
		tp.OriginPrice = -1
		tp.MarketPrice = -1
		tp.IsSale = 0
		tp.IsVisible = 0
		tp.IsRefund = 10

		So(tx, ShouldNotBeNil)
		err = d.TxBatchAddTicketPrice(c, tx, []*model.TicketPrice{tp})
		So(err, ShouldBeNil)

		tps, err := d.TxGetPriceSymbols(c, tx, 1633)
		So(len(tps)-cnt, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})

	Convey("TxGetTicketPrice", t, func() {
		So(tx, ShouldNotBeNil)
		tp, err := d.TxGetTicketPrice(c, tx, genID)
		So(tp.ScreenID, ShouldEqual, 1633)
		So(err, ShouldBeNil)
	})

	tx.Rollback()

	Convey("TxGetTicketPrice", t, func() {
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)
		_, err := d.TxGetTicketPrice(c, tx, genID)
		So(err, ShouldBeError)
	})
}
