package service

import (
	"context"
	"go-common/app/admin/main/aegis/model/net"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicecreateTaskVerify(t *testing.T) {
	var (
		rid    = int64(1)
		flowID = int64(1)
		bizid  = int64(1)
	)
	convey.Convey("prepareBeforeTrigger", t, func(ctx convey.C) {
		err := s.prepareBeforeTrigger(context.TODO(), []int64{rid}, flowID, bizid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceShowTransition(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("ShowTransition", t, func(ctx convey.C) {
		r, err := s.ShowTransition(c, id)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceGetTransitionList(t *testing.T) {
	var (
		c  = context.TODO()
		pm = &net.ListNetElementParam{NetID: 1}
	)
	convey.Convey("GetTransitionList", t, func(ctx convey.C) {
		result, err := s.GetTransitionList(c, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestServicecheckTransitionUnique(t *testing.T) {
	var (
		netID = int64(1)
		name  = ""
	)
	convey.Convey("checkTransitionUnique", t, func(ctx convey.C) {
		err, msg := s.checkTransitionUnique(cntx, netID, name)
		ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
