package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_BuildsByAppName(t *testing.T) {
	svr := svr(t)
	Convey("should app by name", t, func() {
		res, err := svr.Builds(2888, "main.common-arch.msm-service", "dev", "sh001")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestService_BuildByID(t *testing.T) {
	svr := svr(t)
	Convey("should app by name", t, func() {
		res, err := svr.Build(1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
