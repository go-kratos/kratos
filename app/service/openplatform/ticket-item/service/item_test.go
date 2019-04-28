package service

import (
	"context"
	"testing"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s   *ItemService
	ctx = context.TODO()
)

// TestInfo
func TestInfo(t *testing.T) {
	Convey("TestInfo", t, func() {
		res, err := s.Info(ctx, &item.InfoRequest{ID: model.DataID})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("TestInfo", t, func() {
		res, err := s.Info(ctx, &item.InfoRequest{ID: model.NoDataID})
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeNil)
	})
}

// Test_Cards
func TestCards(t *testing.T) {
	Convey("TestCards", t, func() {
		res, err := s.Cards(ctx, &item.CardsRequest{IDs: model.DataIDs})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("TestCards", t, func() {
		res, err := s.Cards(ctx, &item.CardsRequest{IDs: model.NoDataIDs})
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeNil)
	})
}

//
func TestBillInfo(t *testing.T) {
	Convey("TestBillInfo", t, func() {
		res, err := s.BillInfo(ctx, &item.BillRequest{
			IDs:   model.DataIDs,
			ScIDs: model.DataSIDs,
			TkIDs: model.DataTIDs,
		})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("TestBillInfo", t, func() {
		res, err := s.BillInfo(ctx, &item.BillRequest{
			IDs:   model.NoDataIDs,
			ScIDs: model.NoDataSIDs,
			TkIDs: model.NoDataTIDs,
		})
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestWish(t *testing.T) {
	Convey("insert wish", t, func() {
		res, err := s.Wish(ctx, &item.WishRequest{
			ItemID: 1,
			MID:    1,
		})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
