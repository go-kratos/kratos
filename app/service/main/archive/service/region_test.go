package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RegionTopCount(t *testing.T) {
	Convey("RegionTopCount", t, func() {
		_, err := s.RegionTopCount(context.TODO(), []int16{1})
		So(err, ShouldBeNil)
	})
}

func Test_DelRegionArc(t *testing.T) {
	Convey("DelRegionArc", t, func() {
		err := s.DelRegionArc(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_AddRegionArc(t *testing.T) {
	Convey("AddRegionArc", t, func() {
		err := s.AddRegionArc(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}
