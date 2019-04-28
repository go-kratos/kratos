package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDao_WebTopPhoto(t *testing.T) {
	convey.Convey("test get web top photo", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.webTopPhotoURL).Reply(200).JSON(`{"code": 0, "s_img":"test_url", "l_img":"test_url"}`)
		mid := int64(282994)
		data, err := d.WebTopPhoto(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}

func TestDao_TopPhoto(t *testing.T) {
	convey.Convey("test get top photo", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.topPhotoURL).Reply(200).JSON(`{"code": 0,"message": "ok","data": {"image_url":"test_url"}}`)
		mid := int64(908085)
		vmid := int64(908085)
		platform := "ios"
		device := "" // pad
		data, err := d.TopPhoto(context.Background(), mid, vmid, platform, device)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}
