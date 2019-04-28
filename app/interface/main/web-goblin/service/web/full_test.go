package web

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Full(t *testing.T) {
	Convey("test Full", t, WithService(func(s *Service) {
		var (
			pn int64 = 1
			ps int64 = 10
		)
		res, err := s.FullShort(context.Background(), pn, ps, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
