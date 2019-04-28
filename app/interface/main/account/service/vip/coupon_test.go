package vip

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceCouponsForPanel
func TestServiceCouponsForPanel(t *testing.T) {
	Convey("TestServiceCouponsForPanel", t, func() {
		res, err := s.CouponsForPanel(context.TODO(), int64(1), int64(96), "pc")
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}
