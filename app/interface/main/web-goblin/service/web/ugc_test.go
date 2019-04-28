package web

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_UgcFull(t *testing.T) {
	Convey("ugc full", t, WithService(func(s *Service) {
		var (
			c  = context.Background()
			pn = int64(1)
			ps = int64(10)
		)
		res, err := s.UgcFull(c, pn, ps, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_UgcIncre(t *testing.T) {
	Convey("pgc incre", t, WithService(func(s *Service) {
		var (
			c     = context.Background()
			pn    = 1
			ps    = 10
			start = int64(1505876448)
			end   = int64(1505876450)
		)
		res, err := s.UgcIncre(c, pn, ps, start, end, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
