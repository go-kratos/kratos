package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RegionTopArcs3(t *testing.T) {
	Convey("RegionTopArcs3", t, func() {
		_, err := s.RegionTopArcs3(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_RegionAllArcs3(t *testing.T) {
	Convey("RegionAllArcs3", t, func() {
		_, err := s.RegionAllArcs3(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_RegionArcs3(t *testing.T) {
	Convey("RegionArcs3", t, func() {
		_, _, err := s.RegionArcs3(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_RegionOriginArcs3(t *testing.T) {
	Convey("RegionOriginArcs3", t, func() {
		_, _, err := s.RegionOriginArcs3(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
	})
}
