package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Push(t *testing.T) {
	svr := svr(t)
	Convey("should push", t, func() {
		err := svr.Push(context.Background(), 2888, "dev", "sh001", "server-1", 1)
		So(err, ShouldBeNil)
	})
}

func TestService_SetToken(t *testing.T) {
	svr := svr(t)
	Convey("should set token", t, func() {
		err := svr.SetToken(context.Background(), 2888, "dev", "sh001", "84c0c277f13111e79d54522233017188")
		So(err, ShouldBeNil)
	})
}

func TestService_Hosts(t *testing.T) {
	svr := svr(t)
	Convey("should get hosts", t, func() {
		_, err := svr.Hosts(context.Background(), 2888, "main.common-arch.msm-service", "dev", "sh001")
		So(err, ShouldBeNil)
	})
}

func TestService_ClearHost(t *testing.T) {
	svr := svr(t)
	Convey("should clear host", t, func() {
		err := svr.ClearHost(context.Background(), 2888, "dev", "sh001")
		So(err, ShouldBeNil)
	})
}
