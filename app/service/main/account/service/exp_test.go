package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAddExp(t *testing.T) {
	convey.Convey("AddExp", t, func() {
		err := s.AddExp(context.TODO(), 1, 10, "other", "other", "test")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestAddMoral(t *testing.T) {
	convey.Convey("AddMoral", t, func() {
		err := s.AddMoral(context.TODO(), 1, -10, "other", "other", "test")
		convey.So(err, convey.ShouldBeNil)
	})
}
