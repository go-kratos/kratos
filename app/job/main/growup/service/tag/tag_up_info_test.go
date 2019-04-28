package tag

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetTagUpInfoMID(t *testing.T) {
	Convey("GetTagUpInfoMID", t, WithService(func(s *Service) {
		_, err := s.GetTagUpInfoMID(context.Background(), []int64{})
		So(err, ShouldNotBeNil)
	}))

	Convey("GetTagUpInfoMID", t, WithService(func(s *Service) {
		_, err := s.GetTagUpInfoMID(context.Background(), []int64{1})
		So(err, ShouldBeNil)
	}))
}
