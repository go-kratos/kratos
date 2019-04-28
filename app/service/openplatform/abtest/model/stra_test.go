package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var stra = []Stra{
	Stra{Precision: 100, Ratio: []int{10, 90}},
	Stra{Precision: 100, Ratio: []int{10, 9}},
}

func TestCheck(t *testing.T) {
	Convey("TestCheck: ", t, func() {
		var checks = []bool{true, false}
		for i, s := range stra {
			got := s.Check()
			So(got, ShouldEqual, checks[i])
		}
	})
}

func TestVersion(t *testing.T) {
	testCase := map[int]int{9: 0, 20: 1}
	s := stra[0]
	Convey("TestVersion: ", t, func() {
		for j, k := range testCase {
			got, _ := s.Version(j)
			So(got, ShouldEqual, k)

		}

		_, err := s.Version(101)
		So(err, ShouldNotBeNil)

		_, err = stra[1].Version(101)
		So(err, ShouldNotBeNil)
	})
}
