package service

import (
	"context"
	"go-common/app/admin/main/esports/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_EditTeam(t *testing.T) {
	Convey("test edit team", t, WithService(func(s *Service) {
		arg := &model.Team{
			ID: 3,
		}
		gids := []int64{3}
		err := s.EditTeam(context.Background(), arg, gids)
		So(err, ShouldBeNil)
	}))
}
