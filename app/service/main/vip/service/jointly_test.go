package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestServiceJointly
func TestServiceJointly(t *testing.T) {
	Convey(" TestServiceJointly ", t, func() {
		res, err := s.Jointly(c)
		t.Logf("res data(%+v)", res)
		So(err, ShouldBeNil)
	})
}
