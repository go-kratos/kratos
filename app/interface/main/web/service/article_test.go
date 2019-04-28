package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/interface/main/web/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ArticleList(t *testing.T) {
	Convey("test article list", t, WithService(func(s *Service) {
		aids := []int64{}
		rid := 2
		sort := 1
		mid := 27515256
		res, err := s.ArticleList(context.Background(), int64(rid), int64(mid), sort, 1, 20, aids)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		Println(len(res))
		for _, v := range res {
			Printf("%d,%+v", v.ID, v)
		}
	}))
}

func TestService_ArticleUpList(t *testing.T) {
	Convey("test article up list", t, WithService(func(s *Service) {
		res, err := s.ArticleUpList(context.Background(), 90085)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, conf.Conf.Rule.ArtUpListCnt)
		if bs, err := json.Marshal(res); err != nil {
			Printf("%+v", err)
		} else {
			Println(string(bs))
		}
	}))
}

func TestService_Categories(t *testing.T) {
	Convey("test categories", t, WithService(func(s *Service) {
		res, err := s.Categories(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		for _, v := range *res {
			Printf("%+v", v)
		}
	}))
}
