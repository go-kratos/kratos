package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestAddAppeal add appeal testing.
func TestAddAppeal(t *testing.T) {
	Convey("TestAddAppeal", t, func() {
		err := s.AddAppeal(context.TODO(), 1, 111, 2222, "测试")
		So(err, ShouldBeNil)
	})
}

// TestAppealState appealstate testing.
func TestAppealState(t *testing.T) {
	Convey("TestAddAppeal", t, func() {
		state, err := s.AppealState(context.TODO(), 1, 111)
		So(err, ShouldBeNil)
		So(state, ShouldEqual, true)
	})
}
