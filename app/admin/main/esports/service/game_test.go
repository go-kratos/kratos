package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/admin/main/esports/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_GameList(t *testing.T) {
	Convey("test game list", t, WithService(func(s *Service) {
		oids := []int64{21, 16}
		data, err := s.gameList(context.Background(), model.TypeTeam, oids)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(data)
		Printf(string(bs))
	}))
}
