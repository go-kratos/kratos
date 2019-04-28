package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceconvertModel(t *testing.T) {
	convey.Convey("convertModel", t, func() {
		res := s.convertModel(nil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceconvertExtraModel(t *testing.T) {
	convey.Convey("convertExtraModel", t, func() {
		res := s.convertExtraModel(nil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
