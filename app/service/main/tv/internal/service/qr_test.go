package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceMakePayParam(t *testing.T) {
	convey.Convey("MakePayParam", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(27515308)
			pid    = int32(1)
			buyNum = int32(1)
			guid   = "bv23456789"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p := s.MakePayParam(c, mid, pid, buyNum, guid, "master")
			ctx.Convey("Then p should not be nil.", func(ctx convey.C) {
				ctx.So(p, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceCreateQr(t *testing.T) {
	convey.Convey("CreateQr", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(27515308)
			pid    = int32(1)
			buyNum = int32(1)
			guid   = "bv23456789"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			qr, err := s.CreateQr(c, mid, pid, buyNum, guid, "master")
			ctx.Convey("Then err should be nil.qr should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(qr, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceCreateGuestQr(t *testing.T) {
	convey.Convey("CreateGuestQr", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			pid    = int32(1)
			buyNum = int32(1)
			guid   = "4567890"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			qr, err := s.CreateGuestQr(c, pid, buyNum, guid, "master")
			ctx.Convey("Then err should be nil.qr should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(qr, convey.ShouldNotBeNil)
			})
		})
	})
}
