package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServicePointRule(t *testing.T) {
	Convey(" PointRule test ", t, func() {
		_, err := s.PointRule(c)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceBuyVipWithPoint
func TestServiceBuyVipWithPoint(t *testing.T) {
	Convey(" BuyVipWithPoint test ", t, func() {
		err := s.BuyVipWithPoint(c, int64(2), 1)
		So(err, ShouldBeNil)
	})
}
