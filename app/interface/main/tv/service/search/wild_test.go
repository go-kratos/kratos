package search

import (
	"context"
	"testing"

	mdlSearch "go-common/app/interface/main/tv/model/search"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_UserSearch(t *testing.T) {
	Convey("test user search", t, WithService(func(s *Service) {
		arg := &mdlSearch.UserSearch{
			Keyword:    "lex",
			Build:      "111",
			SearchType: "all",
			Page:       1,
			Pagesize:   20,
		}
		res, err := s.UserSearch(context.Background(), arg)
		So(err, ShouldBeNil)
		So(len(res), ShouldNotBeEmpty)
	}))
}

func TestService_SearchAll(t *testing.T) {
	Convey("test search all", t, WithService(func(s *Service) {
		arg := &mdlSearch.UserSearch{
			Keyword:    "工作细胞",
			Build:      "111",
			SearchType: "bili_user",
			Page:       1,
			Pagesize:   20,
		}
		res, err := s.SearchAll(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestService_PgcSearch(t *testing.T) {
	Convey("test pgc search", t, WithService(func(s *Service) {
		arg := &mdlSearch.UserSearch{
			Keyword:    "工作细胞",
			Build:      "111",
			SearchType: "all",
			Page:       1,
			Pagesize:   20,
		}
		res, err := s.PgcSearch(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
