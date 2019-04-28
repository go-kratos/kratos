package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpperReommend(t *testing.T) {
	Convey("UpperReommend", t, func() {
		_, err := s.UpperReommend(context.TODO(), 1)
		So(err, ShouldNotBeNil)
	})
}

func Test_UpperPassed3(t *testing.T) {
	Convey("UpperPassed3", t, func() {
		_, err := s.UpperPassed3(context.TODO(), 1684013, 1, 20)
		So(err, ShouldNotBeNil)
	})
}

func Test_UppersPassed3(t *testing.T) {
	Convey("UppersPassed3", t, func() {
		_, err := s.UppersPassed3(context.TODO(), []int64{1684013}, 1, 20)
		So(err, ShouldBeNil)
	})
}
