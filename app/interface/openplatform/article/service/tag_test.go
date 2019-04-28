package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// func Test_Tag(t *testing.T) {
// 	Convey("BindTags", t, WithService(func(s *Service) {
// 		err := s.BindTags(context.TODO(), 1, 1, []string{"tag1", "tag2", "tag3"}, "")
// 		So(err, ShouldBeNil)
// 	}))
// 	Convey("Tags", t, WithService(func(s *Service) {
// 		aids := []int64{1}
// 		res, err := s.Tags(context.TODO(), aids)
// 		So(err, ShouldBeNil)
// 		So(res, ShouldNotBeEmpty)
// 		len := len(res[1])
// 		t.Logf("result: %+v, tags length: %d", res, len)
// 		So(len, ShouldEqual, 3)
// 	}))
// 	Convey("BindTags 2", t, WithService(func(s *Service) {
// 		err := s.BindTags(context.TODO(), 1, 1, []string{"tag4", "tag5"}, "")
// 		So(err, ShouldBeNil)
// 	}))
// 	Convey("Tags 2", t, WithService(func(s *Service) {
// 		aids := []int64{1}
// 		res, err := s.Tags(context.TODO(), aids)
// 		So(err, ShouldBeNil)
// 		So(res, ShouldNotBeEmpty)
// 		len := len(res[1])
// 		t.Logf("result: %+v, tags length: %d", res, len)
// 		So(len, ShouldEqual, 2)
// 	}))
// }

func Test_mergeActivityTags(t *testing.T) {
	Convey("no act tag add tags", t, func() {
		tags := mergeActivityTags([]string{"tag1", "tag2", "tag3"}, []string{"act1", "act2"})
		So(tags, ShouldResemble, []string{"tag1", "tag2", "tag3", "act1", "act2"})
	})
	Convey("1 tag add tags", t, func() {
		tags := mergeActivityTags([]string{"tag1", "tag2", "tag3", "act1"}, []string{"act2"})
		So(tags, ShouldResemble, []string{"tag1", "tag2", "tag3", "act1", "act2"})
	})
	Convey("already has tags ", t, func() {
		tags := mergeActivityTags([]string{"tag1", "tag2", "tag3", "act1", "act2"}, []string{"act1", "act2"})
		So(tags, ShouldResemble, []string{"tag1", "tag2", "tag3", "act1", "act2"})
	})
	Convey("act tags blank ", t, func() {
		tags := mergeActivityTags([]string{"tag1", "tag2", "tag3"}, []string{})
		So(tags, ShouldResemble, []string{"tag1", "tag2", "tag3"})
	})
}
