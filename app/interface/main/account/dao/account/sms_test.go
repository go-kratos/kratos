package account

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAccountvcodeKey(t *testing.T) {
	var (
		mid    = int64(0)
		mobile = ""
	)
	convey.Convey("vcodeKey", t, func(ctx convey.C) {
		p1 := vcodeKey(mid, mobile)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountSendMobileVerify(t *testing.T) {
	var (
		vcode   = int64(1234)
		country = int64(86)
		mobile  = "13488888888"
		ip      = ""
	)
	convey.Convey("SendMobileVerify", t, func(ctx convey.C) {
		err := d.SendMobileVerify(context.Background(), vcode, country, mobile, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAccountOffiVerifyParam(t *testing.T) {
	var (
		vcode = int64(0)
	)
	convey.Convey("OffiVerifyParam", t, func(ctx convey.C) {
		p1, err := OffiVerifyParam(vcode)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountGenVerifyCode(t *testing.T) {
	var (
		mid    = int64(0)
		mobile = ""
	)
	convey.Convey("GenVerifyCode", t, func(ctx convey.C) {
		p1, err := d.GenVerifyCode(context.Background(), mid, mobile)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountGetVerifyCode(t *testing.T) {
	var (
		mid    = int64(0)
		mobile = ""
	)
	convey.Convey("GetVerifyCode", t, func(ctx convey.C) {
		p1, err := d.GetVerifyCode(context.Background(), mid, mobile)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountDelVerifyCode(t *testing.T) {
	var (
		mid    = int64(0)
		mobile = ""
	)
	convey.Convey("DelVerifyCode", t, func(ctx convey.C) {
		err := d.DelVerifyCode(context.Background(), mid, mobile)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
