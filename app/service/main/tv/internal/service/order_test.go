package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicemakeOrderNo(t *testing.T) {
	convey.Convey("makeOrderNo", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := s.makeOrderNo()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

//
//func TestServicecreateOrder(t *testing.T) {
//	convey.Convey("createOrder", t, func(ctx convey.C) {
//		var (
//			c           = context.Background()
//			mid         = int64(0)
//			token       = ""
//			platform    = int8(0)
//			paymentType = ""
//			clientIp    = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			pi, err := s.createOrder(c, mid, token, platform, paymentType, clientIp)
//			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(pi, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceCreateOrder(t *testing.T) {
//	convey.Convey("CreateOrder", t, func(ctx convey.C) {
//		var (
//			c           = context.Background()
//			token       = ""
//			platform    = int8(0)
//			paymentType = ""
//			clientIp    = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			pi, err := s.CreateOrder(c, token, platform, paymentType, clientIp)
//			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(pi, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceCreateGuestOrder(t *testing.T) {
//	convey.Convey("CreateGuestOrder", t, func(ctx convey.C) {
//		var (
//			c           = context.Background()
//			mid         = int64(0)
//			token       = ""
//			platform    = int8(0)
//			paymentType = ""
//			clientIp    = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			pi, err := s.CreateGuestOrder(c, mid, token, platform, paymentType, clientIp)
//			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(pi, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
