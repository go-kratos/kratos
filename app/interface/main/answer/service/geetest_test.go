package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicePreProcess(t *testing.T) {
	convey.Convey("PreProcess", t, func() {
		res, err := s.preProcess(context.Background(), 0, "", "", 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServicevalidate(t *testing.T) {
	convey.Convey("validate", t, func() {
		stat := s.validate(context.Background(), "", "", "", "", "", 0, 0)
		convey.So(stat, convey.ShouldNotBeNil)
	})
}

func TestServicefailbackValidate(t *testing.T) {
	convey.Convey("failbackValidate", t, func() {
		p1 := s.failbackValidate(context.Background(), "", "", "")
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestServicedecodeResponse(t *testing.T) {
	convey.Convey("decodeResponse", t, func() {
		res := s.decodeResponse("", "")
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServicedecodeRandBase(t *testing.T) {
	convey.Convey("decodeRandBase", t, func() {
		p1 := s.decodeRandBase("")
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestServicemd5Encode(t *testing.T) {
	convey.Convey("md5Encode", t, func() {
		p1 := s.md5Encode([]byte{})
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestServicevalidateFailImage(t *testing.T) {
	convey.Convey("validateFailImage", t, func() {
		p1 := s.validateFailImage(0, 0, 0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}
