package service

import (
	"context"
	"go-common/app/service/main/archive/api"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Stat3(t *testing.T) {
	Convey("Stat3", t, func() {
		_, err := s.Stat3(context.TODO(), 14761597)
		So(err, ShouldBeNil)
	})
}

func Test_Click3(t *testing.T) {
	Convey("Click3", t, func() {
		_, err := s.Click3(context.TODO(), 14761597)
		So(err, ShouldBeNil)
	})
}

func Test_Stats3(t *testing.T) {
	Convey("Stats3", t, func() {
		_, err := s.Stats3(context.TODO(), []int64{14761597})
		So(err, ShouldBeNil)
	})
}

func TestSetStat(t *testing.T) {
	Convey("SetStat", t, func() {
		err := s.SetStat(context.TODO(), &api.Stat{Aid: 1})
		So(err, ShouldBeNil)
	})
}
