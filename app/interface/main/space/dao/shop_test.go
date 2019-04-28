package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDao_ShopInfo(t *testing.T) {
	convey.Convey("test get shop info", t, func(ctx convey.C) {
		mid := int64(27515399)
		defer gock.OffAll()
		httpMock("GET", d.shopURL).Reply(200).JSON(`{"code":0,"data":{"shop":{"id":111}}}`)
		data, err := d.ShopInfo(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}

func TestDao_ShopLink(t *testing.T) {
	convey.Convey("test shop link", t, func(ctx convey.C) {
		mid := int64(27515399)
		plat := 1
		data, err := d.ShopLink(context.Background(), mid, plat)
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldNotBeNil)
		convey.Printf("%+v", data)
	})
}
