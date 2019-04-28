package dao

import (
	"context"
	"go-common/app/service/openplatform/ticket-item/model"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDaoAreaSeats
func TestDaoAreaSeats(t *testing.T) {
	var id int64

	Convey("TxGetAreaSeats", t, func() {
		once.Do(startService)
		c := context.TODO()

		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		area := &model.Area{
			AID:   "test" + strconv.FormatInt(time.Now().Unix(), 10),
			Name:  "可删",
			Place: 132,
		}
		err := d.TxAddArea(c, tx, area)
		So(err, ShouldBeNil)
		So(area.ID, ShouldNotEqual, 0)
		id = area.ID

		res, err := d.TxGetAreaSeats(c, tx, id)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 0)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)
	})
	var seatsIDs []int64
	Convey("TxBatchAddAreaSeats", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		seatsArr := []*model.AreaSeats{
			{
				X:       1,
				Y:       2,
				Label:   "3",
				Bgcolor: "#FFFFFF",
				Area:    id,
				Dstatus: 0,
			},
			{
				X:       2,
				Y:       2,
				Label:   "4",
				Bgcolor: "#FFFFFF",
				Area:    id,
				Dstatus: 0,
			},
		}
		err := d.TxBatchAddAreaSeats(c, tx, seatsArr)
		So(err, ShouldBeNil)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		tx, _ = d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		res, err := d.TxGetAreaSeats(c, tx, id)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)

		seatsIDs = []int64{res[0].ID, res[1].ID}
		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)
	})

	Convey("TxBatchDeleteAreaSeats", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		err := d.TxBatchDeleteAreaSeats(c, tx, id)
		So(err, ShouldBeNil)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		tx, _ = d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		res, err := d.TxGetAreaSeats(c, tx, id)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)
		So(res[0].Dstatus, ShouldEqual, 1)
		So(res[1].Dstatus, ShouldEqual, 1)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)
	})

	Convey("TxBatchRecoverAreaSeats", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		err := d.TxBatchRecoverAreaSeats(c, tx, seatsIDs)
		So(err, ShouldBeNil)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		tx, _ = d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		res, err := d.TxGetAreaSeats(c, tx, id)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)
		So(res[0].Dstatus, ShouldEqual, 0)
		So(res[1].Dstatus, ShouldEqual, 0)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)
	})
}

// TestDaoAreaSeats
func TestDaoAreaSeatmap(t *testing.T) {
	once.Do(startService)
	c := context.TODO()
	tx, _ := d.BeginTran(c)

	Convey("TxRawAreaSeatmap", t, func() {
		So(tx, ShouldNotBeNil)

		asm, err := d.TxRawAreaSeatmap(c, tx, 239)
		So(asm.ID, ShouldEqual, 239)
		So(asm.SeatMap, ShouldNotBeBlank)
		So(err, ShouldBeNil)
	})

	Convey("TxRawAreaSeatmap", t, func() {
		So(tx, ShouldNotBeNil)

		asm := &model.AreaSeatmap{
			ID:      239,
			SeatMap: "[\"abc\"]",
		}
		err := d.TxSaveAreaSeatmap(c, tx, asm)
		So(err, ShouldBeNil)

		asm, err = d.TxRawAreaSeatmap(c, tx, 239)
		So(asm.ID, ShouldEqual, 239)
		So(asm.SeatMap, ShouldEqual, "[\"abc\"]")
		So(err, ShouldBeNil)
	})

	tx.Rollback()
}

func TestSeatSet(t *testing.T) {
	once.Do(startService)
	c := context.TODO()
	tx, _ := d.BeginTran(c)

	Convey("TxGetSeatChart", t, func() {
		So(tx, ShouldNotBeNil)
		sc, err := d.TxGetSeatChart(c, tx, 1633, 239)
		So(sc.ID, ShouldBeZeroValue)
		So(sc.SeatChart, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("TxGetSeatCharts", t, func() {
		So(tx, ShouldNotBeNil)
		scs, err := d.TxGetSeatCharts(c, tx, 1633, []int64{239})
		So(len(scs), ShouldBeZeroValue)
		So(err, ShouldBeNil)
	})

	Convey("TxAddSeatChart", t, func() {
		So(tx, ShouldNotBeNil)
		err := d.TxAddSeatChart(c, tx, &model.SeatSet{
			AreaID:   239,
			ScreenID: 1633,
		})
		So(err, ShouldBeNil)
	})

	var id int64
	Convey("TxGetSeatCharts", t, func() {
		So(tx, ShouldNotBeNil)
		scs, err := d.TxGetSeatCharts(c, tx, 1633, []int64{239})
		So(len(scs), ShouldEqual, 1)
		So(scs[0].SeatChart, ShouldBeEmpty)
		So(err, ShouldBeNil)
		id = scs[0].ID
	})

	Convey("TxUpdateSeatChart", t, func() {
		So(tx, ShouldNotBeNil)
		err := d.TxUpdateSeatChart(c, tx, id, "[\"_______aaaaaaaaaaaaa\",\"_______aaaaaaaaaaaaa\",\"______aaaaaaaaaaaaaa\",\"______aaaaaaaaaaaaaa\",\"_____aaaaaaaaaaaaaaa\",\"_____aaaaaaaaaaaaaaa\",\"____aaaaaaaaaaaaaaaa\",\"____aaaaaaaaaaaaaaaa\",\"___aaaaaaaaaaaaaaaaa\",\"___aaaaaaaaaaaaaaaaa\",\"__aaaaaaaaaaaaaaaaaa\",\"__aaaaaaaaaaaaaaaaaa\",\"aaaaaaaaaaaaaaaaaaaa\",\"aaaaaaaaaaaaaaaaaaaa\",\"___aaaaaaaaaaaaaa___\",\"_aaaaaaaaaaaaaaaaaa_\"]")
		So(err, ShouldBeNil)
	})

	Convey("TxGetSeatChart", t, func() {
		So(tx, ShouldNotBeNil)
		sc, err := d.TxGetSeatChart(c, tx, 1633, 239)
		So(sc.ID, ShouldEqual, id)
		So(sc.SeatChart, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})

	tx.Rollback()
}

// TestSeatOrder
func TestSeatOrder(t *testing.T) {
	once.Do(startService)
	c := context.TODO()
	tx, _ := d.BeginTran(c)

	Convey("TxBatchAddSeatOrder", t, func() {
		So(tx, ShouldNotBeNil)
		err := d.TxBatchAddSeatOrder(c, tx, []*model.SeatOrder{
			{
				AreaID:   239,
				ScreenID: 1633,
				Row:      1,
				Col:      1,
				PriceID:  36093,
				Price:    123456,
			}, {
				AreaID:   239,
				ScreenID: 1633,
				Row:      2,
				Col:      2,
				PriceID:  36093,
				Price:    123456,
			}, {
				AreaID:   239,
				ScreenID: 1633,
				Row:      3,
				Col:      3,
				PriceID:  36093,
				Price:    123456,
			},
		})
		So(err, ShouldBeNil)
	})

	Convey("TxGetUnsaleableSeatOrders", t, func() {
		res, err := d.TxGetUnsaleableSeatOrders(c, tx, 1633, 239)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 0)
	})

	var ids [3]int64
	Convey("TxGetSaleableSeatOrders", t, func() {
		res, err := d.TxGetSaleableSeatOrders(c, tx, 1633, 36093)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 3)
		for i, so := range res {
			ids[i] = so.ID
		}
	})

	Convey("TxBatchDeleteSeatOrder", t, func() {
		err := d.TxBatchDeleteSeatOrder(c, tx, ids[:])
		So(err, ShouldBeNil)
	})

	tx.Commit()
	tx, _ = d.BeginTran(c)
	Convey("TxGetSaleableSeatOrders", t, func() {
		res, err := d.TxGetSaleableSeatOrders(c, tx, 1633, 36093)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 0)
	})
	tx.Rollback()
}
