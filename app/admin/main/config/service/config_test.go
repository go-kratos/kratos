package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ConfigsByIDs(t *testing.T) {
	svr := svr(t)
	Convey("should configs by ids", t, func() {
		res, err := svr.ConfigsByIDs([]int64{1, 2})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})

}

func TestService_UpdateConfState(t *testing.T) {
	svr := svr(t)
	Convey("should update state", t, func() {
		err := svr.UpdateConfState(2)
		So(err, ShouldBeNil)
	})
}

func TestService_ConfigsByBuildID(t *testing.T) {
	svr := svr(t)
	Convey("should get config by build id", t, func() {
		res, err := svr.ConfigsByBuildID(1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestService_ConfigsByAppName(t *testing.T) {
	svr := svr(t)
	Convey("should get configs by app name", t, func() {
		res, err := svr.ConfigsByAppName("main.common-arch.msm-service", "dev", "shd", 2888, 0)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
func TestService_Configs(t *testing.T) {
	svr := svr(t)
	Convey("should configs", t, func() {
		res, err := svr.Configs("main.account.open-svr-mng", "fat1", "shd", 0, 2888)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestService_Diff(t *testing.T) {
	svr := svr(t)
	Convey("should configs", t, func() {
		res, err := svr.Diff(1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
func TestService_Value(t *testing.T) {
	svr := svr(t)
	Convey("should configs", t, func() {
		res, err := svr.Value(1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
