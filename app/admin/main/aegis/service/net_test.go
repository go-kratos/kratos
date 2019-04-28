package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/net"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceGetNetList(t *testing.T) {
	var (
		c  = context.TODO()
		pm = &net.ListNetParam{}
	)
	convey.Convey("GetNetList", t, func(ctx convey.C) {
		result, err := s.GetNetList(c, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceShowNet(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("ShowNet", t, func(ctx convey.C) {
		r, err := s.ShowNet(c, id)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestServicenetCheckUnique(t *testing.T) {
	var (
		chName = ""
	)
	convey.Convey("netCheckUnique", t, func(ctx convey.C) {
		err, msg := s.netCheckUnique(cntx, chName)
		ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
