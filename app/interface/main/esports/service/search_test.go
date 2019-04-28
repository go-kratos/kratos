package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/esports/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Search(t *testing.T) {
	Convey("test service Search", t, WithService(func(s *Service) {
		arg := &model.ParamSearch{
			Keyword: "dota",
			Pn:      1,
			Ps:      30,
		}
		res, err := s.Search(context.Background(), 0, arg, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
