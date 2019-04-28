package dao

import (
	"context"
	"go-common/app/service/openplatform/ticket-item/model"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestRawVenues
func TestRawVenues(t *testing.T) {
	Convey("RawVenues", t, func() {
		once.Do(startService)
		res, err := d.RawVenues(context.TODO(), []int64{77, 78})
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)
		So(res[77].Name, ShouldEqual, "海尔测试-勿动")
	})
}

// TestDaoVenue
func TestDaoVenue(t *testing.T) {
	var id int64
	Convey("AddVenue", t, func() {
		once.Do(startService)
		venue := model.Venue{
			Name:          "可删",
			Province:      100000,
			City:          100000,
			District:      100000,
			Status:        0,
			AddressDetail: "地址详细",
			Traffic:       "交通",
		}
		err := d.AddVenue(context.TODO(), &venue)
		So(err, ShouldBeNil)
		So(venue.ID, ShouldNotEqual, 0)
		id = venue.ID
	})

	Convey("UpdateVenue", t, func() {
		once.Do(startService)
		venue := model.Venue{
			ID:            id,
			Name:          "可删",
			Province:      100000,
			City:          100000,
			District:      100010,
			Status:        1,
			AddressDetail: "地址详细",
			Traffic:       "交通",
		}
		err := d.UpdateVenue(context.TODO(), &venue)
		So(err, ShouldBeNil)
		So(venue.Status, ShouldEqual, 1)
	})

	var placeNum int32

	Convey("TxIncPlaceNum", t, func() {
		once.Do(startService)
		c := context.TODO()
		res, err := d.RawVenues(c, []int64{id})
		placeNum = res[id].PlaceNum
		So(err, ShouldBeNil)

		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)
		err = d.TxIncPlaceNum(c, tx, id)
		So(err, ShouldBeNil)
		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		res, err = d.RawVenues(c, []int64{id})
		So(err, ShouldBeNil)
		So(res[id].PlaceNum, ShouldEqual, placeNum+1)
	})

	Convey("TxDecPlaceNum", t, func() {
		once.Do(startService)
		c := context.TODO()

		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)
		err := d.TxDecPlaceNum(c, tx, id)
		So(err, ShouldBeNil)
		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		res, _ := d.RawVenues(c, []int64{id})
		So(res[id].PlaceNum, ShouldEqual, placeNum)
	})
}

// TestRawPlace
func TestRawPlace(t *testing.T) {
	Convey("RawPlace", t, func() {
		once.Do(startService)
		res, err := d.RawPlace(context.TODO(), 132)
		So(err, ShouldBeNil)
		So(res.Name, ShouldEqual, "海尔测试-勿动")
	})
}

// TestTxRawPlace
func TestTxRawPlace(t *testing.T) {
	Convey("TxRawPlace", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		res, err := d.TxRawPlace(c, tx, 132)
		So(err, ShouldBeNil)
		So(res.Name, ShouldEqual, "海尔测试-勿动")

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		So(res.Name, ShouldEqual, "海尔测试-勿动")
	})
}

// TestDaoPlace
func TestDaoPlace(t *testing.T) {
	var id int64
	Convey("TxAddPlace", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		place := &model.Place{
			Name:    "可删",
			BasePic: "url",
			Status:  1,
			Venue:   81,
			DWidth:  800,
			DHeight: 800,
		}
		err := d.TxAddPlace(c, tx, place)
		So(err, ShouldBeNil)
		So(place.ID, ShouldNotEqual, 0)
		id = place.ID

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		place, err = d.RawPlace(c, id)
		So(err, ShouldBeNil)
		So(place.ID, ShouldNotEqual, 0)
	})

	Convey("TxUpdatePlace", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		place := &model.Place{
			ID:      id,
			Name:    "可删",
			BasePic: "url2",
			Status:  1,
			Venue:   81,
			DWidth:  800,
			DHeight: 800,
		}
		err := d.TxUpdatePlace(c, tx, place)
		So(err, ShouldBeNil)
		So(place.BasePic, ShouldEqual, "url2")

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		place, err = d.RawPlace(c, id)
		So(err, ShouldBeNil)
		So(place.BasePic, ShouldEqual, "url2")
	})

	var placePolygon *model.PlacePolygon

	Convey("TxAddOrUpdateAreaPolygon", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		var co = "292 103,443 113,414 263,230 247"
		err := d.TxAddOrUpdateAreaPolygon(c, tx, id, 239, &co)
		So(err, ShouldBeNil)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		placePolygon, err = d.RawPlacePolygon(c, id)
		So(err, ShouldBeNil)
		So(placePolygon.Coordinate, ShouldEqual, "{\"239\":\"292 103,443 113,414 263,230 247\"}")
	})

	Convey("TxDelAreaPolygon", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		err := d.TxDelAreaPolygon(c, tx, id, 239)
		So(err, ShouldBeNil)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		placePolygon, err = d.RawPlacePolygon(c, id)
		So(err, ShouldBeNil)
		So(placePolygon.Coordinate, ShouldEqual, "{}")
	})
}

// TestRawArea
func TestRawArea(t *testing.T) {
	Convey("RawArea", t, func() {
		once.Do(startService)
		res, err := d.RawArea(context.TODO(), 239)
		So(err, ShouldBeNil)
		So(res.Name, ShouldEqual, "海尔测试-勿动")
	})
}

func TestTxRawArea(t *testing.T) {
	Convey("RawArea", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		res, err := d.TxRawArea(c, tx, 239)
		So(err, ShouldBeNil)
		So(res.Name, ShouldEqual, "海尔测试-勿动")

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		res, err = d.RawArea(context.TODO(), 239)
		So(err, ShouldBeNil)
		So(res.Name, ShouldEqual, "海尔测试-勿动")
	})
}

// TestDaoArea
func TestDaoArea(t *testing.T) {
	var id int64
	Convey("TxAddArea", t, func() {
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

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		area, err = d.RawArea(c, id)
		So(err, ShouldBeNil)
		So(area.ID, ShouldNotEqual, 0)
	})

	Convey("TxUpdateArea", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		area := &model.Area{
			ID:    id,
			AID:   "test" + strconv.FormatInt(time.Now().Unix(), 10),
			Name:  "可删-1",
			Place: 132,
		}
		err := d.TxUpdateArea(c, tx, area)
		So(err, ShouldBeNil)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		area, err = d.RawArea(c, id)
		So(err, ShouldBeNil)
		So(area.Name, ShouldEqual, "可删-1")
		So(area.Place, ShouldEqual, 132)
	})

	Convey("TxUpdateSeat", t, func() {
		once.Do(startService)
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		So(tx, ShouldNotBeNil)

		area := &model.Area{
			ID:            id,
			SeatsNum:      10,
			Width:         1,
			Height:        2,
			DeletedStatus: 0,
			ColStart:      3,
			ColType:       4,
			ColDirection:  1,
			RowList:       "[\"1\",\"2\",\"3\",\"4\",\"5\",\"6\",\"7\",\"8\",\"9\",\"10\",\"11\",\"12\",\"13\",\"14\",\"15\",\"16\",\"17\",\"18\",\"19\",\"20\"]",
			SeatStart:     "",
		}
		err := d.TxUpdateSeat(c, tx, area)
		So(err, ShouldBeNil)
		So(area.SeatsNum, ShouldEqual, 10)

		area, err = d.RawArea(c, id)
		So(err, ShouldBeNil)
		So(area.SeatsNum, ShouldEqual, 0)

		err = d.CommitTran(c, tx)
		So(err, ShouldBeNil)

		area, err = d.RawArea(c, id)
		So(err, ShouldBeNil)
		So(area.SeatsNum, ShouldEqual, 10)
	})

}
