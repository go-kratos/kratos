package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Article(t *testing.T) {
	Convey("article list test", t, WithService(func(s *Service) {
		mid := int64(442549)
		pn := 1
		ps := 10
		sort := model.ArticleSortType["publish_time"]
		res, err := s.Article(context.Background(), mid, pn, ps, sort)
		So(err, ShouldBeNil)
		if res != nil && len(res.Articles) > 0 {
			Print(len(res.Articles), res.Count)
			for _, v := range res.Articles {
				Printf("%+v", v)
			}
		}
	}))
}
