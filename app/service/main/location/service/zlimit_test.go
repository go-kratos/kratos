package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PgcZone(t *testing.T) {
	Convey("pgc zone in", t, WithService(func(s *Service) {
		_, err := s.PgcZone(context.Background(), []int64{2, 124425})
		So(err, ShouldBeNil)
	}))
}

func Test_Auth(t *testing.T) {
	Convey("get archive auth", t, WithService(func(s *Service) {
		_, err := s.Auth(context.Background(), 740955, 13414510, "211.139.80.6", "64.233.173.24")
		So(err, ShouldBeNil)
	}))
}

func Test_AuthGID(t *testing.T) {
	Convey("get group auth", t, WithService(func(s *Service) {
		res := s.AuthGID(context.Background(), 1195, 13414510, "211.139.80.6", "64.233.173.24")
		So(res, ShouldNotBeNil)
	}))
}
