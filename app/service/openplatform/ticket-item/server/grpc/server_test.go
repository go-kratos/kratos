package server

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/app/service/openplatform/ticket-item/service"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s   *server
	ctx = context.TODO()
)

func init() {
	dir, _ := filepath.Abs("../../cmd/item.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	s = &server{
		is: service.New(conf.Conf),
	}
}

// Test_Info
func TestInfo(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.Info(ctx, &item.InfoRequest{ID: model.DataID})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("no data", t, func() {
		res, err := s.Info(ctx, &item.InfoRequest{ID: model.NoDataID})
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// Test_Cards
func TestCards(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.Cards(ctx, &item.CardsRequest{IDs: model.DataIDs})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("no data", t, func() {
		res, err := s.Cards(ctx, &item.CardsRequest{IDs: model.NoDataIDs})
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeNil)
	})
}

// TestBillInfo
func TestBillInfo(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.BillInfo(ctx, &item.BillRequest{
			IDs:   model.DataIDs,
			ScIDs: model.DataSIDs,
			TkIDs: model.DataTIDs,
		})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("no data", t, func() {
		res, err := s.BillInfo(ctx, &item.BillRequest{
			IDs:   model.NoDataIDs,
			ScIDs: model.NoDataSIDs,
			TkIDs: model.NoDataTIDs,
		})
		So(err, ShouldNotBeNil)
		So(res, ShouldNotBeNil)
	})
}

// TestWish
func TestWish(t *testing.T) {
	Convey("add wish", t, func() {
		res, err := s.Wish(ctx, &item.WishRequest{
			ItemID: 1,
			MID:    1,
			Face:   "",
		})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

// TestWish
func TestFav(t *testing.T) {
	Convey("fav add", t, func() {
		res, err := s.Fav(ctx, &item.FavRequest{
			ItemID: 1,
			MID:    1,
			Type:   1,
			Status: 1,
		})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})

	Convey("fav del", t, func() {
		res, err := s.Fav(ctx, &item.FavRequest{
			ItemID: 1,
			MID:    1,
			Type:   1,
			Status: 0,
		})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
