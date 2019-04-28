package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestAllGroup
func TestAllGroup(t *testing.T) {
	Convey("TestAllGroup ", t, func() {
		card, err := s.AllGroup(c, 1)
		t.Logf("v(%v)", card)
		So(err, ShouldBeNil)
	})
}
