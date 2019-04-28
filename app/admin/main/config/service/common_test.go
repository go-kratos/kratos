package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ComConfigsByTagID(t *testing.T) {
	svr := svr(t)
	Convey("should app by name", t, func() {
		res, err := svr.ConfigsByTagID(1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestService_ComConfigsByTeam(t *testing.T) {
	svr := svr(t)
	Convey("should get common config by team", t, func() {
		res, err := svr.ComConfigsByTeam("main.common-arch", "dev", "shd", 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
