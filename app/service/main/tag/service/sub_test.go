package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Sub(t *testing.T) {
	Convey("AddSub", t, WithService(func(s *Service) {
		err := s.AddSub(context.Background(), mid, tids, ip)
		So(err, ShouldBeNil)
	}))
	Convey("CancelSub", t, WithService(func(s *Service) {
		err := s.CancelSub(context.Background(), mid, tid, ip)
		So(err, ShouldBeNil)
	}))
	Convey("SubTags", t, WithService(func(s *Service) {
		_, _, err := s.SubTags(context.Background(), mid, pn, ps, order)
		So(err, ShouldBeNil)
	}))
}
