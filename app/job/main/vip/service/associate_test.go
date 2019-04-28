package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//go test  -test.v -test.run  TestEleEompensateJob
func TestEleEompensateJob(t *testing.T) {
	Convey("TestEleEompensateJob ", t, func() {
		s.eleEompensateJob()
	})
}

//go test  -test.v -test.run  TestEleGrantCompensate
func TestEleGrantCompensate(t *testing.T) {
	Convey("TestEleGrantCompensate ", t, func() {
		err := s.EleGrantCompensate(c)
		So(err, ShouldNotBeNil)
	})
}
