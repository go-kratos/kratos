package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/esports/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Calendar(t *testing.T) {
	Convey("test service calendar", t, WithService(func(s *Service) {
		res, err := s.Calendar(context.Background(), &model.ParamFilter{Stime: "2018-07-27", Etime: "2018-08-02"})
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_FilterMatch(t *testing.T) {
	Convey("test service filterMatch", t, WithService(func(s *Service) {
		res, err := s.FilterMatch(context.Background(), &model.ParamFilter{Mid: 0})
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
func TestService_FilterVideo(t *testing.T) {
	Convey("test service filterVideo", t, WithService(func(s *Service) {
		res, err := s.FilterVideo(context.Background(), &model.ParamFilter{Mid: 0})
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_ListVideo(t *testing.T) {
	Convey("test service listVideo", t, WithService(func(s *Service) {
		arg := &model.ParamVideo{
			Mid:  int64(0),
			Gid:  int64(0),
			Tid:  int64(0),
			Year: int64(2018),
			Tag:  int64(1),
			Pn:   1,
			Ps:   30,
		}
		res, total, err := s.ListVideo(context.Background(), arg)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		println(total)
	}))
}

func TestService_ListContest(t *testing.T) {
	Convey("test service listContest", t, WithService(func(s *Service) {
		arg := &model.ParamContest{
			Mid:    int64(0),
			GState: "0,3,4",
			Pn:     1,
			Ps:     10,
		}
		mid := int64(12309)
		res, total, err := s.ListContest(context.Background(), mid, arg)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		println(total)
	}))
}

func TestService_Season(t *testing.T) {
	Convey("test service Season", t, WithService(func(s *Service) {
		arg := &model.ParamSeason{
			Pn: 1,
			Ps: 5,
		}
		res, count, err := s.Season(context.Background(), arg)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
