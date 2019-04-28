package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetNotices(t *testing.T) {
	var (
		typ           = 0
		platform      = 1
		offset, limit = 0, 10
	)
	Convey("GetNotices", t, WithService(func(s *Service) {
		res, _, err := s.GetNotices(context.Background(), typ, platform, offset, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
