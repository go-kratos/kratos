package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_LoadRes(t *testing.T) {
	Convey("load resource", t, WithService(func(s *Service) {
		err := s.loadRes()
		So(err, ShouldBeNil)
	}))
}

func Test_ResourceAll(t *testing.T) {
	Convey("get all resources ", t, WithService(func(s *Service) {
		res := s.ResourceAll(context.Background())
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_AssignmentAll(t *testing.T) {
	Convey("get all assignment", t, WithService(func(s *Service) {
		res := s.AssignmentAll(context.Background())
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_Resource(t *testing.T) {
	Convey("get resource by id", t, WithService(func(s *Service) {
		res := s.Resource(context.Background(), 2329)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_Resources(t *testing.T) {
	Convey("get resource by ids", t, WithService(func(s *Service) {
		res := s.Resources(context.Background(), []int{467, 2329})
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_DefBanner(t *testing.T) {
	Convey("get default banner", t, WithService(func(s *Service) {
		res := s.DefBanner(context.Background())
		So(res, ShouldNotBeNil)
	}))
}

func Test_IndexIcon(t *testing.T) {
	Convey("get index icon", t, WithService(func(s *Service) {
		res := s.IndexIcon(context.Background())
		So(res, ShouldNotBeEmpty)
	}))
}
