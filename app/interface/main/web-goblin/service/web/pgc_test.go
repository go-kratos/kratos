package web

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_PgcFull(t *testing.T) {
	Convey("pgc full", t, WithService(func(s *Service) {
		var (
			c      = context.Background()
			tp     = int(2)
			pn     = int64(1)
			ps     = int64(10)
			source = "youku"
		)
		res, err := s.PgcFull(c, tp, pn, ps, source)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_PgcIncre(t *testing.T) {
	Convey("pgc incre", t, WithService(func(s *Service) {
		var (
			c      = context.Background()
			tp     = int(2)
			pn     = int64(1)
			ps     = int64(10)
			start  = int64(1505876448)
			end    = int64(1505876450)
			source = "youku"
		)
		res, err := s.PgcIncre(c, tp, pn, ps, start, end, source)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
