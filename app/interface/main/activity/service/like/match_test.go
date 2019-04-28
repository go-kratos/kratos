package like

import (
	"context"
	"testing"

	"go-common/app/interface/main/activity/model/like"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Match(t *testing.T) {
	Convey("test service Match", t, WithService(func(s *Service) {
		sid := int64(1)
		res, err := s.Match(context.Background(), sid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_AddGuess(t *testing.T) {
	Convey("test service  AddGuess", t, WithService(func(s *Service) {
		mid := int64(111)
		objID := int64(1)
		result := int64(2)
		stake := int64(3)
		lastID, err := s.AddGuess(context.Background(), mid, &like.ParamAddGuess{ObjID: objID, Result: result, Stake: stake})
		So(err, ShouldBeNil)
		So(lastID, ShouldBeGreaterThan, 0)
	}))
}

func TestService_ListGuess(t *testing.T) {
	Convey("test service  match ListGuess", t, WithService(func(s *Service) {
		sid := int64(111)
		mid := int64(111)
		guess, err := s.ListGuess(context.Background(), sid, mid)
		So(err, ShouldBeNil)
		So(len(guess), ShouldBeGreaterThan, 0)
	}))
}

func TestService_Guess(t *testing.T) {
	Convey("test service  match Guess", t, WithService(func(s *Service) {
		mid := int64(111)
		sid := int64(1)
		guess, err := s.Guess(context.Background(), mid, &like.ParamSid{Sid: sid})
		So(err, ShouldBeNil)
		So(guess, ShouldNotBeNil)
	}))
}

func TestService_ClearCache(t *testing.T) {
	Convey("test service ClearCache", t, WithService(func(s *Service) {
		msg := `{"action":"update","table":"act_matchs_object","old":{"name":0},"new":{"id":2,"sid":12,"match_id":2}}`
		err := s.ClearCache(context.Background(), msg)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddFollow(t *testing.T) {
	Convey("test service AddFollow", t, WithService(func(s *Service) {
		mid := int64(111)
		team := []string{"11", "22", "33"}
		err := s.AddFollow(context.Background(), mid, team)
		So(err, ShouldBeNil)
	}))
}

func TestService_Follow(t *testing.T) {
	Convey("test service Follow", t, WithService(func(s *Service) {
		mid := int64(111)
		teams, err := s.Follow(context.Background(), mid)
		So(err, ShouldBeNil)
		So(len(teams), ShouldBeGreaterThan, 0)
	}))
}

func TestService_ObjectsUnStart(t *testing.T) {
	Convey("test service ObjectsUnStart", t, WithService(func(s *Service) {
		mid := int64(111)
		sid := int64(1)
		objs, total, err := s.ObjectsUnStart(context.Background(), mid, &like.ParamObject{Sid: sid, Pn: 1, Ps: 4})
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(objs), ShouldBeGreaterThan, 0)
	}))
}
