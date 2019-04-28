package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpperCount(t *testing.T) {
	Convey("UpperCount", t, func() {
		_, err := s.UpperCount(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_UppersAidPubTime(t *testing.T) {
	Convey("UppersAidPubTime", t, func() {
		_, err := s.UppersAidPubTime(context.TODO(), []int64{1684013}, 1, 10)
		So(err, ShouldBeNil)
	})
}

func Test_AddUpperPassedCache(t *testing.T) {
	Convey("AddUpperPassedCache", t, func() {
		err := s.AddUpperPassedCache(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_DelUpperPassedCache(t *testing.T) {
	Convey("DelUpperPassedCache", t, func() {
		err := s.DelUpperPassedCache(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_UpperCache(t *testing.T) {
	Convey("UpperCache", t, func() {
		err := s.UpperCache(context.TODO(), 1684913, "updateByAdmin")
		So(err, ShouldBeNil)
	})
}

func Test_UppersCache(t *testing.T) {
	Convey("UppersCache", t, func() {
		uc, err := s.UppersCount(context.TODO(), []int64{168403, 15555180})
		So(err, ShouldBeNil)
		Println(uc)
	})
}
