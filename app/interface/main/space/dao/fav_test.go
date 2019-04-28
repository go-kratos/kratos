package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDao_AlbumFavCount(t *testing.T) {
	convey.Convey("test album fav count", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.favAlbumURL).Reply(200).JSON(`{"code": 0, "data":{"pageinfo":{"count":1}}}`)
		mid := int64(88895029)
		favType := 2
		data, err := d.LiveFavCount(context.Background(), mid, favType)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}

func TestDao_MovieFavCount(t *testing.T) {
	convey.Convey("test movie fav count", t, func(ctx convey.C) {
		mid := int64(88895029)
		favType := 2
		data, err := d.MovieFavCount(context.Background(), mid, favType)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}

func TestDao_FavFolder(t *testing.T) {
	convey.Convey("test fav folder", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.favFolderURL).Reply(200).JSON(`{"code": 0, "data":[]}`)
		mid := int64(88895029)
		vmid := int64(0)
		data, err := d.FavFolder(context.Background(), mid, vmid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}

func TestDao_FavArchive(t *testing.T) {
	convey.Convey("test fav archive", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.favArcURL).Reply(200).JSON(`{"code": 0}`)
		mid := int64(88895029)
		arg := &model.FavArcArg{}
		data, err := d.FavArchive(context.Background(), mid, arg)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}
