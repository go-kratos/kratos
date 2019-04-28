package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestService_Ping(t *testing.T) {
	Convey("normal", t, testWith(func() {
		err := s.Ping(ctx)
		So(err, ShouldBeNil)
	}))
}
