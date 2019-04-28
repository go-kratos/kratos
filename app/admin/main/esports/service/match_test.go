package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/esports/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddMatch(t *testing.T) {
	Convey("test add match", t, WithService(func(s *Service) {
		gids := []int64{1}
		err := s.AddMatch(context.Background(), &model.Match{Title: "match"}, gids)
		So(err, ShouldBeNil)
	}))
}

func TestService_ForbidMatch(t *testing.T) {
	Convey("test forbid match", t, WithService(func(s *Service) {
		mid := int64(1)
		state := 0
		err := s.ForbidMatch(context.Background(), mid, state)
		So(err, ShouldBeNil)
	}))
}
