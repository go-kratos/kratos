package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_handlerPayOrder(t *testing.T) {
	Convey("handler pay order", t, func() {
		s.HandlerPayOrder()
	})
}

func Test_autoRenews(t *testing.T) {
	Convey("autorenews ", t, func() {
		err := s.autoRenews()
		So(err, ShouldBeNil)
	})
}
