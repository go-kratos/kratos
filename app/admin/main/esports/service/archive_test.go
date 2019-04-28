package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/esports/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_arcRelationChanges(t *testing.T) {
	Convey("test arc relation", t, WithService(func(s *Service) {
		arg := &model.ArcImportParam{
			Aid:      2147483647,
			Gids:     []int64{2, 3, 5},
			MatchIDs: []int64{91},
			TeamIDs:  []int64{2, 3, 4},
			TagIDs:   []int64{2, 4},
			Years:    []int64{1986, 2003},
		}
		data, err := s.arcRelationChanges(context.Background(), arg, _typeEdit)
		So(err, ShouldBeNil)
		for _, v := range data.AddTeams {
			Printf("Add %+v \n", v)
		}
		for _, v := range data.UpAddTeams {
			Printf("upAdd %+v \n", v)
		}
		for _, v := range data.UpDelTeams {
			Printf("upDel %+v \n", v)
		}
	}))
}

func TestService_BatchAddArc(t *testing.T) {
	Convey("test batch add arc", t, WithService(func(s *Service) {
		arg := &model.ArcAddParam{
			Aids:     []int64{44444444, 55555555},
			Gids:     []int64{6},
			MatchIDs: []int64{31},
			TeamIDs:  []int64{3, 2},
			TagIDs:   []int64{3},
			Years:    []int64{2017, 2018},
		}
		err := s.BatchAddArc(context.Background(), arg)
		So(err, ShouldBeNil)
	}))
}
