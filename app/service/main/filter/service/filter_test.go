package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_CoverStart(t *testing.T) {
	Convey("Test_CoverStart", t, func() {
		sbh := &sortBlackHits{}
		sbh.Add([]int{4, 6})
		out, err := sbh.CoverByStart("test23test")
		So(out, ShouldEqual, "test**test")
		So(err, ShouldBeNil)
	})
}

func Test_CoverStart_Fail(t *testing.T) {
	Convey("Test_CoverStart", t, func() {
		sbh := &sortBlackHits{}
		sbh.Add([]int{4, 6})
		out, err := sbh.CoverByStart("test23test")
		So(out, ShouldNotEqual, "test23test")
		So(err, ShouldBeNil)
	})
}
