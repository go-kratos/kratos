package operation

import (
	"context"
	"go-common/library/ecode"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_WebOperations(t *testing.T) {
	Convey("WebOperations", t, WithService(func(s *Service) {
		res, err := s.WebOperations(context.TODO())
		So(err, ShouldEqual, ecode.NothingFound)
		So(res, ShouldBeNil)
	}))
	Convey("Ping", t, WithService(func(s *Service) {
		err := s.Ping(context.TODO())
		So(err, ShouldEqual, ecode.NothingFound)
	}))
}

func Test_AppBanner(t *testing.T) {
	Convey("AppBanner", t, WithService(func(s *Service) {
		bns, cbns, err := s.AppBanner(context.TODO())
		So(bns, ShouldBeNil)
		So(cbns, ShouldBeNil)
		So(err, ShouldNotBeNil)
	}))
}

func Test_CreatorOperationList(t *testing.T) {
	Convey("CreatorOperationList", t, WithService(func(s *Service) {
		bns, err := s.CreatorOperationList(context.TODO(), 1, 10)
		So(bns, ShouldBeNil)
		So(err, ShouldNotBeNil)
	}))
}
