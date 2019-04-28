package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoOpenCode(t *testing.T) {
	var (
		c           = context.TODO()
		mid         = int64(0)
		batchCodeID = int64(0)
		unit        = int32(0)
		remark      = ""
		code        = ""
	)
	convey.Convey("OpenCode", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _openCode).Reply(200).JSON(`{"code":0}`)
		data, err := d.OpenCode(c, mid, batchCodeID, unit, remark, code)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("data should not be nil", func(ctx convey.C) {
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetOpenInfo(t *testing.T) {
	var (
		c    = context.TODO()
		code = ""
	)
	convey.Convey("GetOpenInfo", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", _openCode).Reply(200).JSON(`{"code":0}`)
		data, err := d.GetOpenInfo(c, code)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("data should not be nil", func(ctx convey.C) {
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetToken(t *testing.T) {
	var (
		c   = context.TODO()
		bid = "abc"
		ip  = ""
	)
	convey.Convey("GetToken", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.c.Property.APICoURL+_token).Reply(200).JSON(`{"code":0}`)
		_, err := d.GetToken(c, bid, ip)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoVerify(t *testing.T) {
	var (
		c     = context.TODO()
		code  = "abc"
		token = "abc"
		ip    = ""
	)
	convey.Convey("Verify", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _verify).Reply(200).JSON(`{"code":0}`)
		data, err := d.Verify(c, code, token, ip)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("data should not be nil", func(ctx convey.C) {
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetPassportDetail(t *testing.T) {
	convey.Convey("passport", t, func() {
		var mid int64 = 27515586
		defer gock.OffAll()
		httpMock("GET", _passportDetail).Reply(200).JSON(`{"code":0}`)
		_, err := d.GetPassportDetail(context.TODO(), mid)
		convey.So(err, convey.ShouldBeNil)
	})
}
