package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpArticleMetas(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		res, err := s.UpArticleMetas(context.TODO(), dataMID, 1, 20, 0)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithService(func(s *Service) {
		res, err := s.UpArticleMetas(context.TODO(), dataMID, 20000, 20, 0)
		So(err, ShouldBeNil)
		So(res.Articles, ShouldBeEmpty)
	}))
}

func Test_UpsArticles(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		res, err := s.UpsArticleMetas(context.TODO(), []int64{dataMID}, 1, 20)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("get count", t, WithService(func(s *Service) {
		res, err := s.UpperArtsCount(context.TODO(), dataMID)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	}))
	Convey("no data", t, WithService(func(s *Service) {
		res, err := s.UpsArticleMetas(context.TODO(), []int64{dataID}, 20000, 20)
		So(err, ShouldBeNil)
		So(res[dataID], ShouldBeEmpty)
	}))
}
