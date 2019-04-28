package account

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

const (
	mid = int64(2089809)
	ip  = "127.0.0.1"
)

func TestPhoneEmail(t *testing.T) {
	var (
		c  = context.TODO()
		ck = "iamck"
	)
	convey.Convey("PhoneEmail", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.passURI).Reply(200).JSON(`{"code":20007}`)
		ret, err := d.PhoneEmail(c, ck, ip)
		ctx.Convey("Then err should be nil.ret should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(ret, convey.ShouldBeNil)
		})
	})
}

func TestAccountPic(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Pic", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.picUpInfoURL).Reply(200).JSON(`{"code":0,"data":{"has_doc":100}}`)
		has, err := d.Pic(c, mid, ip)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(has, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountBlink(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Blink", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.blinkUpInfoURL).Reply(200).JSON(`{"code":0,"data":{"has":100}}`)
		has, err := d.Blink(c, mid, ip)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(has, convey.ShouldNotBeNil)
		})
	})
}
