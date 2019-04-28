package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Regions(t *testing.T) {
	Convey("test service regions", t, WithService(func(s *Service) {
		res, err := s.Regions(context.Background())
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
