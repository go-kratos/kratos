package dao

import (
	"context"
	"go-common/library/ecode"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDao_Live(t *testing.T) {
	convey.Convey("test live", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.liveURL).Reply(200).JSON(`{"code": 0}`)
		mid := int64(28272030)
		platform := "ios"
		data, err := d.Live(context.Background(), mid, platform)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}

func TestDao_LiveMetal(t *testing.T) {
	convey.Convey("test live metal", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.liveMetalURL).Reply(200).JSON(`{"code": 510002}`)
		mid := int64(28272030)
		data, err := d.LiveMetal(context.Background(), mid)
		convey.So(err, convey.ShouldEqual, ecode.Int(510002))
		convey.Printf("%v", data)
	})
}
