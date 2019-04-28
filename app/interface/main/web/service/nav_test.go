package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Nav(t *testing.T) {
	Convey("test nav Nav", t, WithService(func(s *Service) {
		var (
			mid int64
			ck  = ""
		)
		res, err := s.Nav(context.Background(), mid, ck)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
